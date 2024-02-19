package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/services/events/processImages/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func handler(ctx context.Context, s3Event events.S3Event) {
	service := service.New(log, s3client)

	for _, record := range s3Event.Records {
		log.Info("Got the s3 event",
			"event", record.EventName,
		)
		service.Log.Info("Got the event")
		service.ProcessImagesS3(ctx, record.S3.Object.Key, record.S3.Bucket.Name)
	}
}

func main() {
	lambda.Start(handler)
}
