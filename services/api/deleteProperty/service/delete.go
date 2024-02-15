package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type DeleteService struct {
	Log            *slog.Logger
	PropertiesRepo *repository.Properties
	SQSClient      *sqs.Client
}

func NewDeleteService(
	log *slog.Logger,
	propertiesRepo *repository.Properties,
	sqsClient *sqs.Client,
) DeleteService {
	return DeleteService{
		Log:            log,
		PropertiesRepo: propertiesRepo,
		SQSClient:      sqsClient,
	}
}

func (s DeleteService) DeleteProperty(ctx context.Context, id, humanReadableId string) error {
	deletQueueUrl := os.Getenv(env.SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL)

	liteProperty, err := s.PropertiesRepo.GetById(ctx, "", humanReadableId)
	if err != nil {
		s.Log.Error(
			"error trying to get the property data before deletion",
			"error", err,
		)
		return err
	}
	// Dispatch the delete images SQS event and then delete the property
	imagesToDeleteList := strings.Split(liteProperty.Images, ";")
	sqsDeleteImagesObject := utils.SQSDeleteImages{
		Images: imagesToDeleteList,
	}
	sqsDeleteImagesJson, err := json.Marshal(sqsDeleteImagesObject)
	s.Log.Info("images to delete dispatched to sqs queue",
		"list", imagesToDeleteList,
		"jsoned", string(sqsDeleteImagesJson),
	)

	if err != nil {
		s.Log.Error(
			"error trying to marshal the sqs delete images",
			"error", err,
		)
		return err
	}

	_, err = s.SQSClient.SendMessage(
		ctx,
		&sqs.SendMessageInput{
			QueueUrl:    &deletQueueUrl,
			MessageBody: aws.String(string(sqsDeleteImagesJson)),
		},
	)

	if err != nil {
		s.Log.Error(
			"error trying to dispatch the sqs delete images message",
			"error", err,
		)
		return err
	}

	err = s.PropertiesRepo.DeleteByHumanReadableId(ctx, humanReadableId)
	if err != nil {
		s.Log.Error("error when trying to delete the property", "error", err)
		return fmt.Errorf("error trying to delete the property with human readable id %w", err)
	}

	// ❗TODO dispatch the SQS request to delete processed images

	return nil
}
