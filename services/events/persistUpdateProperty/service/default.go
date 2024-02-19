package service

import (
	"context"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

// TODO handle the property creation
// handle the checking if the user id exists before persisting the property

// TODO make sure all fields are mapped from repository model
// TODO move this struct to the handler
type PUService struct {
	Log                *slog.Logger
	PropertyRepository *repository.Properties
	SQSClient          *sqs.Client
}

func NewPUService(
	log *slog.Logger,
	pr *repository.Properties,
	sqsClient *sqs.Client,
) PUService {
	return PUService{
		Log:                log.WithGroup("Submit Service"),
		PropertyRepository: pr,
		SQSClient:          sqsClient,
	}
}

// 💩 TODO the logic to delete images is flawed
// Review the logic
func (ss *PUService) Update(ctx context.Context, sqsBody string, humanReadableId string) error {
	sqsDeleteImagesQueue := os.Getenv(env.SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL)
	var requestProperty utils.RequestProperty

	if err := json.Unmarshal([]byte(sqsBody), &requestProperty); err != nil {
		ss.Log.Error("error decoding the json property for PUT request", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}
	if err := json.Unmarshal([]byte(sqsBody), &requestProperty.Features); err != nil {
		ss.Log.Error("error decoding the json property features for PUT request", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}

	features, err := json.Marshal(requestProperty.Features)
	if err != nil {
		return err
	}

	// get the existing property to append the existing files if new files are added
	liteProperty, err := ss.PropertyRepository.GetById(ctx, "", humanReadableId)
	if err != nil {
		ss.Log.Error("error trying to get the existing proprety", "hrID", humanReadableId, "error", err)
		return fmt.Errorf("error trying to get the existing property to update")
	}

	var finalImages []string = make([]string, 0)
	finalImages = append(finalImages, strings.Split(liteProperty.Images, ";")...)
	// if there are no new images added the ImageNames should be 0
	if len(requestProperty.ImageNames) > 0 {
		finalImages = append(finalImages, utils.AppendMultipleFileExtension(requestProperty.ImageNames, "webp")...)
	}
	var filteredImages = make([]string, 0)
	removeMap := make(map[string]bool, 0)
	for _, img := range requestProperty.DeletedImages {
		removeMap[img] = true
	}

	for _, img := range finalImages {
		if !removeMap[img] { // if one of the final images exists with a true value in removeMap it's not included in filteredImages
			filteredImages = append(filteredImages, img)
		}
	}

	if requestProperty.DeletedImages != nil && len(requestProperty.DeletedImages) > 0 {

		sqsDeleteImagesObject := utils.SQSDeleteImages{
			Images: requestProperty.DeletedImages,
		}

		jsonDeletedImagesList, err := json.Marshal(sqsDeleteImagesObject)
		ss.Log.Info("images to delete dispatched to sqs queue",
			"list", requestProperty.DeletedImages,
			"jsoned", string(jsonDeletedImagesList),
		)

		if err != nil {
			ss.Log.Error(
				"error trying to marshal the sqs delete images",
				"error", err,
			)
			return err
		}
		_, err = ss.SQSClient.SendMessage(
			ctx,
			&sqs.SendMessageInput{
				QueueUrl:    &sqsDeleteImagesQueue,
				MessageBody: aws.String(string(jsonDeletedImagesList)),
			},
		)

		if err != nil {
			ss.Log.Error(
				"error trying to dispatch the sqs delete images message",
				"error", err,
			)
			return err
		}
	}

	// TODO handle the case where no new images are uploaded
	// OR when new images are uploaded and the old ones are not deleted
	if err := ss.PropertyRepository.UpdateProperty(ctx, lite.UpdatePropertyFieldsParams{
		Humanreadableid:     humanReadableId,
		Title:               requestProperty.Title,
		Images:              strings.Join(filteredImages, ";"),
		Thumbnail:           filteredImages[0],
		IsFeatured:          utils.BoolToInt(requestProperty.IsFeatured),
		PropertyTransaction: requestProperty.PropertyTransaction,
		PropertyDescription: requestProperty.PropertyDescription,
		PropertyType:        requestProperty.PropertyType,
		PropertyAddress:     requestProperty.PropertyAddress,
		PropertySurface:     int64(requestProperty.PropertySurface),
		Price:               int64(requestProperty.Price),
		Features:            string(features),
		UpdatedAt:           time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf("error trying to persist the order with error: %v", err)
	}
	return nil
}

func (ss *PUService) Persist(ctx context.Context, sqsBody string, userId, humanReadableId string) error {
	var propertyId = uuid.New().String()
	var requestProperty utils.RequestProperty
	thumbnail := ""

	if err := json.Unmarshal([]byte(sqsBody), &requestProperty); err != nil {
		ss.Log.Error("error decoding the json property", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}
	if err := json.Unmarshal([]byte(sqsBody), &requestProperty.Features); err != nil {
		ss.Log.Error("error decoding the json property features for PUT request", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}

	// ❗TODO
	// is it needed to marshal the features again?
	features, err := json.Marshal(requestProperty.Features)
	if err != nil {
		return err
	}

	if len(requestProperty.ImageNames) > 0 {
		thumbnail = utils.AppendFileExtension(requestProperty.ImageNames[0], "webp")
	}

	if err := ss.PropertyRepository.Add(ctx, lite.AddPropertyParams{
		ID:              propertyId,
		UserID:          userId,
		Humanreadableid: humanReadableId,
		Title:           requestProperty.Title,
		Images: strings.Join(
			utils.AppendMultipleFileExtension(requestProperty.ImageNames, "webp"),
			";",
		),
		// TODO
		// Until I can find a way to figure out when all the images for a property have successfully processd
		// I'll just use the IsProcessing true flag
		IsProcessing:        0, // It's always going to be 0 until S3 images are processed
		Thumbnail:           thumbnail,
		IsFeatured:          utils.BoolToInt(requestProperty.IsFeatured),
		PropertyTransaction: requestProperty.PropertyTransaction,
		PropertyDescription: requestProperty.PropertyDescription,
		PropertyType:        requestProperty.PropertyType,
		PropertyAddress:     requestProperty.PropertyAddress,
		PropertySurface:     int64(requestProperty.PropertySurface),
		Price:               int64(requestProperty.Price),
		Features:            string(features),
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf("error trying to persist the order with error: %v", err)
	}
	return nil
}
