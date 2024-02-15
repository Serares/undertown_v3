package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	rootUtils "github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func handler(ctx context.Context, sqsDeleteImagesEvent events.SQSEvent) error {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	PROCESSED_IMAGES_BUCKET := os.Getenv(env.PROCESSED_IMAGES_BUCKET)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error(
			"error trying to load the default config",
			"error", err,
		)
	}
	client := s3.NewFromConfig(cfg)

	for _, message := range sqsDeleteImagesEvent.Records {
		// TODO is it better to send an array of strings
		// or send a sqs message for each image?
		var deleteImages rootUtils.SQSDeleteImages
		err := json.Unmarshal([]byte(message.Body), &deleteImages)

		if err != nil {
			log.Error(
				"error trying to unmarshal the sqs body",
				"error", err,
			)
			return err
		}

		for _, imageToDelete := range deleteImages.Images {
			processedImageKey := constants.S3_PROCESSED_IMAGES_PREFIX + "/" + imageToDelete
			_, err := client.DeleteObject(
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
				return err
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
