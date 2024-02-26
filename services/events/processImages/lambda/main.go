package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"runtime"
	"sync"

	"github.com/Serares/undertown_v3/services/events/processImages/service"
	"github.com/Serares/undertown_v3/utils"
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

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	service := service.New(log, s3client)

	for _, record := range sqsEvent.Records {
		log.Info("Got the s3 event",
			"event", record.MessageId,
		)
		var rawImagesMessage utils.SQSImagesMessage

		err := json.Unmarshal([]byte(record.Body), &rawImagesMessage)
		if err != nil {
			log.Error("error unmarshaling the sqs message",
				"error", err,
			)
		}
		var wg sync.WaitGroup
		var errChan = make(chan error, len(rawImagesMessage.Images))
		var s3Errors = make([]error, 0)
		sem := make(chan interface{}, runtime.NumCPU())

		for _, rawImageName := range rawImagesMessage.Images {
			wg.Add(1)
			go func(rawImageName string) {
				defer wg.Done()
				sem <- struct{}{}
				service.ProcessImagesS3(
					ctx,
					rawImageName,
					rawImagesMessage.HumanReadableId,
					errChan,
				)
				<-sem

			}(rawImageName)
		}
		wg.Wait()
		close(errChan)             // meaning no more errors will be sent on the channel/the loop is over
		for err := range errChan { // this is going to run indeffinetly until a close signal is met on the channel
			if err != nil {
				s3Errors = append(s3Errors, err)
			}
		}
		if len(s3Errors) > 0 {
			err := errors.Join(s3Errors...)
			log.Error("Error on s3 processing images", "errors", err)
			return err // this should send the sqs message to dlq
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
