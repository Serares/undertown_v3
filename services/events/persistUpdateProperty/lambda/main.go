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
)

// listening for PIU_QUEUE
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbRepo, err := repository.NewPropertiesRepo()
	ss := service.NewSubmitService(log, dbRepo)

	if err != nil {
		log.Error("error trying to connect to db %v", err)
	}
	// ❗TODO
	// to error handle this use a DLQ
	for _, message := range sqsEvent.Records {
		// of the user id is part of sqs metadata
		// it meas it's a persist message
		// else it's an update message
		if userId, ok := message.MessageAttributes[constants.USER_ID]; ok {
			if humanReadableId, ok := message.MessageAttributes[constants.HUMAN_READABLE_ID_SQS_ATTRIBUTE]; ok {
				err := ss.ProcessPropertyData(
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
				}
			}
		} else if humanReadableId, ok := message.MessageAttributes[constants.HUMAN_READABLE_ID_SQS_ATTRIBUTE]; ok {
			err := ss.ProcessPropertyUpdateData(ctx, message.Body, *humanReadableId.StringValue)
			if err != nil {
				log.Error(
					"error trying to update the property from the sqs message",
					"error", err,
					"hrID", humanReadableId,
				)
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
