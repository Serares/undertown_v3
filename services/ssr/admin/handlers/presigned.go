package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	rootUtils "github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PresignedS3Handler struct {
	Log           *slog.Logger
	PresignClient *s3.PresignClient
}

func NewPresignedS3Handler(
	log *slog.Logger,
	presignClient *s3.PresignClient,
) *PresignedS3Handler {
	return &PresignedS3Handler{
		Log:           log.WithGroup("Presign Handler"),
		PresignClient: presignClient,
	}
}

type PresignedRequest struct {
	FileNames []string `json:"fileNames"`
}

type PresignedResponse struct {
	Response map[string]string `json:"response"` // map the fileName : presignedUrl
}

type PresignedResponseValue struct {
	KeyName      string `json:"KeyName"`
	PresignedUrl string `json:"PresignedUrl"`
}

type presignedResult struct {
	OriginalFileName string
	KeyName          string
	PresignedUrl     string
	Error            error
}

// type PresignFilekey struct {
// 	PresignedUrl string `json:"presignedUrl"`
// 	Filename     string `json:"fileName"`
// }

func (service *PresignedS3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bucketName := os.Getenv(env.RAW_IMAGES_BUCKET)
		var presignReq PresignedRequest

		err := json.NewDecoder(r.Body).Decode(&presignReq)
		if err != nil {
			service.Log.Error("error trying to decode the json request", "err", err)
			rootUtils.ReplyError(
				w,
				r,
				http.StatusInternalServerError,
				"Error decoding the reuquest",
			)
			return
		}
		var fileNameKeyNameMap map[string]string = make(map[string]string)

		for _, fileName := range presignReq.FileNames {
			keyName := rootUtils.ReplaceWhiteSpaceWithUnderscore(
				fmt.Sprintf(
					"%s_%s",
					rootUtils.GenerateStringTimestamp(),
					fileName,
				),
			)
			fileNameKeyNameMap[fileName] = keyName
		}

		var presignedResponse map[string]PresignedResponseValue = make(map[string]PresignedResponseValue)
		var s3PresignWg sync.WaitGroup
		var s3Responses = make(chan presignedResult, len(presignReq.FileNames))
		var s3ResponseErrors = make([]error, 0)

		for fileName, keyName := range fileNameKeyNameMap {
			s3PresignWg.Add(1)
			go func(keyName, fileName string) {
				defer s3PresignWg.Done()
				presignParameters := &s3.PutObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(keyName),
				}

				presignResponse, err := service.PresignClient.PresignPutObject(
					r.Context(),
					presignParameters,
					s3.WithPresignExpires(time.Minute*5),
				)
				s3Responses <- presignedResult{
					OriginalFileName: fileName,
					KeyName:          keyName,
					PresignedUrl:     presignResponse.URL,
					Error:            err,
				}
			}(keyName, fileName)
		}

		go func() {
			s3PresignWg.Wait()
			close(s3Responses)
		}()

		for result := range s3Responses {
			if result.Error == nil {
				presignedResponse[result.OriginalFileName] = PresignedResponseValue{
					KeyName:      result.KeyName,
					PresignedUrl: result.PresignedUrl,
				}
			} else {
				s3ResponseErrors = append(s3ResponseErrors, result.Error)
			}
		}

		if len(s3ResponseErrors) > 0 {
			service.Log.Error(
				"Error requesting s3 presigned url",
				"errors",
				errors.Join(s3ResponseErrors...),
				"success responses",
				presignedResponse,
			)
			rootUtils.ReplyError(
				w,
				r,
				http.StatusInternalServerError,
				"Error on S3 Presign request",
			)
			return
		}

		rootUtils.ReplySuccess(
			w,
			r,
			http.StatusAccepted,
			presignedResponse,
		)
	}
}
