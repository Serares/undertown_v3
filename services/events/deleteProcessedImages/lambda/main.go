package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"sync"

	rootUtils "github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3client *s3.Client
var log *slog.Logger

func init() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error(
			"error trying to load the default config",
			"error", err,
		)
	}
	s3client = s3.NewFromConfig(cfg)
}

func handler(ctx context.Context, sqsDeleteImagesEvent events.SQSEvent) error {
	PROCESSED_IMAGES_BUCKET := os.Getenv(env.PROCESSED_IMAGES_BUCKET)

	for _, message := range sqsDeleteImagesEvent.Records {
		var deleteImages rootUtils.SQSImagesMessage
		err := json.Unmarshal([]byte(message.Body), &deleteImages)

		if err != nil {
			log.Error(
				"error trying to unmarshal the sqs body",
				"error", err,
			)
			return err
		}

		var wg sync.WaitGroup
		var s3errors = make([]error, 0)
		var s3errorsChannel = make(chan error, len(deleteImages.Images))

		for _, imageToDelete := range deleteImages.Images {
			wg.Add(1)
			go func(imageToDelete string) {
				defer wg.Done()
				processedImageKey := constants.S3_PROCESSED_IMAGES_PREFIX + "/" + imageToDelete
				_, err := s3client.DeleteObject(
					ctx,
					&s3.DeleteObjectInput{
						Bucket: aws.String(PROCESSED_IMAGES_BUCKET),
						Key:    aws.String(processedImageKey),
					},
				)
				if err != nil {
					log.Error("error trying to delete the key:",
						"key:", imageToDelete,
					)
					s3errorsChannel <- err
				}
			}(imageToDelete)
		}

		wg.Wait()
		close(s3errorsChannel)
		for err := range s3errorsChannel {
			if err != nil {
				s3errors = append(s3errors, err)
			}
		}
		if len(s3errors) > 0 {
			err := errors.Join(s3errors...)
			log.Error("Error on s3 processing images", "errors", err)
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
