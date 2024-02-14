package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/services/events/processImages/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
	)

	s3Client, err := service.NewS3Client(
		ctx,
	)

	if err != nil {
		log.Error("error trying to initialize the s3 client with config",
			"error", err,
		)
	}

	service := service.New(log, s3Client)

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
