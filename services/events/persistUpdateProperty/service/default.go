package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"strings"
	"sync"
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

func (ss *PUService) Update(ctx context.Context, sqsBody string, humanReadableId string) error {
	sqsDeleteImagesQueue := os.Getenv(env.SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL)
	sqsProcessRawImages := os.Getenv(env.SQS_PROCESS_RAW_IMAGES_QUEUE_URL)

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

	var finalImages = make([]string, 0)
	var hydratedImages = make([]string, 0) //

	finalImages = append(finalImages, strings.Split(liteProperty.Images, ";")...)
	// if there are no new images added the ImageNames should be 0
	if len(requestProperty.ImageNames) > 0 {
		hydratedImages =
			utils.PrependImagesWithHrId(
				utils.ReplaceFileExtensionForList(
					requestProperty.ImageNames, "webp"), humanReadableId)

		finalImages = append(finalImages, hydratedImages...)
	}

	var prefixedDeletedImages = utils.PrependImagesWithHrId(requestProperty.DeletedImages, humanReadableId)
	var filteredImages = make([]string, 0)
	removeMap := make(map[string]bool, 0)
	for _, img := range prefixedDeletedImages {
		removeMap[img] = true
	}

	for _, img := range finalImages {
		if !removeMap[img] { // if one of the final images exists with a true value in removeMap it's not included in filteredImages
			filteredImages = append(filteredImages, img)
		}
	}

	var wgSqs sync.WaitGroup
	sqsErrors := make([]error, 2)
	wgSqs.Add(2)
	go func() {
		defer wgSqs.Done()
		if len(prefixedDeletedImages) > 0 {
			sqsDeleteImagesObject := utils.SQSDeleteImages{
				Images: prefixedDeletedImages,
			}

			jsonDeletedImagesList, err := json.Marshal(sqsDeleteImagesObject)
			ss.Log.Info("images to delete dispatched to sqs queue",
				"list", prefixedDeletedImages,
				"jsoned", string(jsonDeletedImagesList),
			)

			if err != nil {
				ss.Log.Error(
					"error trying to marshal the sqs delete images",
					"error", err,
				)
				sqsErrors[0] = err
				return
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
				sqsErrors[0] = err
				return
			}
		}
	}()
	go func() {
		defer wgSqs.Done()
		if len(requestProperty.ImageNames) > 0 {
			sqsRawImages := utils.SQSProcessRawImages{
				Images:          requestProperty.ImageNames,
				HumanReadableId: humanReadableId,
			}

			jsonRawImages, err := json.Marshal(sqsRawImages)
			if err != nil {
				sqsErrors[1] = err
				return
			}
			// send the raw images to be processed
			_, err = ss.SQSClient.SendMessage(
				ctx,
				&sqs.SendMessageInput{
					QueueUrl:    &sqsProcessRawImages,
					MessageBody: aws.String(string(jsonRawImages)),
				},
			)
			if err != nil {
				sqsErrors[1] = err
				return
			}
		}
	}()
	wgSqs.Wait()

	// ⚠️TODO find a better way of handling this async errors
	err = errors.Join(sqsErrors...)
	if err != nil {
		return err
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

func (ss *PUService) Persist(ctx context.Context, sqsBody string, userId string) error {
	var propertyId = uuid.New().String()
	var requestProperty utils.RequestProperty
	sqsProcessRawImagesUrl := os.Getenv(env.SQS_PROCESS_RAW_IMAGES_QUEUE_URL)

	if err := json.Unmarshal([]byte(sqsBody), &requestProperty); err != nil {
		ss.Log.Error("error decoding the json property", "err", err)
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

	humanReadableId := utils.HumanReadableId(requestProperty.PropertyType)

	// prefix with the hrID because that's how the processed images are stored in S3
	images := utils.PrependImagesWithHrId(requestProperty.ImageNames, humanReadableId)
	s3ImagesPathsProcessed := utils.HydrateImagesNames(utils.ReplaceFileExtensionForList(images, "webp"))

	// dispatch sqs process image sqs message
	rawImagesSqsMessage := utils.SQSProcessRawImages{
		Images:          requestProperty.ImageNames,
		HumanReadableId: humanReadableId,
	}

	jsonRawImagesSqsMessage, err := json.Marshal(rawImagesSqsMessage)
	if err != nil {
		return err
	}

	_, err = ss.SQSClient.SendMessage(
		ctx,
		&sqs.SendMessageInput{
			QueueUrl:    &sqsProcessRawImagesUrl,
			MessageBody: aws.String(string(jsonRawImagesSqsMessage)),
		},
	)
	if err != nil {
		return err
	}
	if err := ss.PropertyRepository.Add(ctx, lite.AddPropertyParams{
		ID:              propertyId,
		UserID:          userId,
		Humanreadableid: humanReadableId,
		Title:           requestProperty.Title,
		Images: strings.Join(
			s3ImagesPathsProcessed,
			";",
		),
		// TODO
		// Until I can find a way to figure out when all the images for a property have successfully processd
		// I'll just use the IsProcessing true flag
		IsProcessing:        0, // It's always going to be 0 until S3 images are processed
		Thumbnail:           s3ImagesPathsProcessed[0],
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
