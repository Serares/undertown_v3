package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/events/persistUpdateProperty/service"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var sqsClient *sqs.Client
var log *slog.Logger
var repo *repository.Properties

// 🚀
// the init function is supposed to be automatically called before the main() function
func init() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error(
			"error trying to load the lambda context",
			"error", err,
		)
	}
	sqsClient = sqs.NewFromConfig(cfg)
	repo, err = repository.NewPropertiesRepo()
	if err != nil {
		log.Error(
			"error trying to initialize the repository",
			"error", err,
		)
	}
}

// ⚠️
// TODO
// move to separating the update and the persist to separate lambdas
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	ss := service.NewPUService(log, repo, sqsClient)

	for _, message := range sqsEvent.Records {
		// ⚠️ TODO
		// don't use message attributes to infer the message type
		// define a sqs message type for update and persist
		// check the type not the attributes
		if userId, ok := message.MessageAttributes[constants.USER_ID]; ok {
			err := ss.Persist(
				ctx,
				message.Body,
				*userId.StringValue,
			)
			if err != nil {
				log.Error(
					"error trying to persist the property from the sqs message",
					"error", err,
				)
				return err
			}
		} else if humanReadableId, ok := message.MessageAttributes[constants.HUMAN_READABLE_ID_SQS_ATTRIBUTE]; ok {
			err := ss.Update(
				ctx,
				message.Body,
				*humanReadableId.StringValue,
			)
			if err != nil {
				log.Error(
					"error trying to update the property from the sqs message",
					"error", err,
					"hrID", humanReadableId,
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
