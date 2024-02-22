package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	FileName string `json:"fileName"`
}

type PresignedResponse struct {
	PresignedUrl string `json:"presignedUrl"`
	KeyName      string `json:"keyName"`
}

func (psh *PresignedS3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bucketName := os.Getenv(env.RAW_IMAGES_BUCKET)
		var presignReq PresignedRequest

		err := json.NewDecoder(r.Body).Decode(&presignReq)
		if err != nil {
			psh.Log.Error("error trying to decode the json request", "err", err)
			rootUtils.ReplyError(
				w,
				r,
				http.StatusInternalServerError,
				"Error decoding the reuquest",
			)
			return
		}

		keyName := rootUtils.ReplaceWhiteSpaceWithUnderscore(
			fmt.Sprintf(
				"%s_%s",
				rootUtils.GenerateStringTimestamp(),
				presignReq.FileName,
			),
		)

		presignParameters := &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(keyName),
		}

		presignResponse, err := psh.PresignClient.PresignPutObject(
			r.Context(),
			presignParameters,
			s3.WithPresignExpires(time.Minute*5),
		)

		if err != nil {
			psh.Log.Error("error generating the presigned url", "error", err)
			rootUtils.ReplyError(
				w,
				r,
				http.StatusInternalServerError,
				"Error on generating the presigned url",
			)
			return
		}

		resp := PresignedResponse{
			PresignedUrl: presignResponse.URL,
			KeyName:      keyName,
		}

		rootUtils.ReplySuccess(
			w,
			r,
			http.StatusAccepted,
			resp,
		)
	}
}
