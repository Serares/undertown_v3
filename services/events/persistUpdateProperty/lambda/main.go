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

// listening for PIU_QUEUE
// sending SQS_
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Error(
			"error trying to load the lambda context",
			"error", err,
		)
		return err
	}
	sqsClient := sqs.NewFromConfig(cfg)

	dbRepo, err := repository.NewPropertiesRepo()
	ss := service.NewPUService(log, dbRepo, sqsClient)

	if err != nil {
		log.Error("error trying to connect to db %v", err)
		return err
	}

	for _, message := range sqsEvent.Records {
		// of the user id is part of sqs metadata
		// it meas it's a persist message
		// else it's an update message
		if userId, ok := message.MessageAttributes[constants.USER_ID]; ok {
			if humanReadableId, ok := message.MessageAttributes[constants.HUMAN_READABLE_ID_SQS_ATTRIBUTE]; ok {
				err := ss.Persist(
					ctx,
					message.Body,
					*userId.StringValue,
					*humanReadableId.StringValue,
				)
				if err != nil {
					log.Error(
						"error trying to persist the property from the sqs message",
						"error", err,
					)
					return err
				}
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
		// get the message attributes
		// the userId
		// the humanReadableId ❗has to be created on admin ssr because the images sent will contain the hrID

	}

	return nil
}

func main() {
	lambda.Start(handler)
}
