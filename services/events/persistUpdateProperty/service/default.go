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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

func (serv *PUService) Update(ctx context.Context, sqsBody string, humanReadableId string) error {
	sqsDeleteImagesQueueUrl := os.Getenv(env.SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL)
	sqsProcessRawImagesUrl := os.Getenv(env.SQS_PROCESS_RAW_IMAGES_QUEUE_URL)

	var requestProperty utils.RequestProperty

	if err := json.Unmarshal([]byte(sqsBody), &requestProperty); err != nil {
		serv.Log.Error("error decoding the json property for PUT request", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}
	if err := json.Unmarshal([]byte(sqsBody), &requestProperty.Features); err != nil {
		serv.Log.Error("error decoding the json property features for PUT request", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}

	features, err := json.Marshal(requestProperty.Features)
	if err != nil {
		return err
	}

	// get the existing property to append the existing files if new files are added
	liteProperty, err := serv.PropertyRepository.GetById(ctx, "", humanReadableId)
	if err != nil {
		serv.Log.Error("error trying to get the existing proprety", "hrID", humanReadableId, "error", err)
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
			sqsDeletedImagesList := serv.SplitImagesIntoSQSImageMessages(prefixedDeletedImages, 2, "")
			sqsMessagesForBatch, err := serv.BatchSQSImageMessagesIntoRequestEntries(sqsDeletedImagesList, 10)
			if err != nil {
				serv.Log.Error(
					"error trying to create the batches of SQS Request Entries",
					"error", err,
				)
			}

			// this is going to wait
			err = serv.DispatchSQSBatches(ctx, sqsMessagesForBatch, sqsDeleteImagesQueueUrl)

			if err != nil {
				serv.Log.Error("error processing the sqs batches",
					"error", err,
				)
				sqsErrors[1] = err
			}
		}
	}()
	go func() {
		defer wgSqs.Done()
		if len(requestProperty.ImageNames) > 0 {
			sqsRawImagesList := serv.SplitImagesIntoSQSImageMessages(requestProperty.ImageNames, 2, humanReadableId)
			sqsMessagesForBatch, err := serv.BatchSQSImageMessagesIntoRequestEntries(sqsRawImagesList, 10)
			if err != nil {
				serv.Log.Error(
					"error trying to create the batches of SQS Request Entries",
					"error", err,
				)
			}

			// this is going to wait
			err = serv.DispatchSQSBatches(ctx, sqsMessagesForBatch, sqsProcessRawImagesUrl)

			if err != nil {
				serv.Log.Error("error processing the sqs batches",
					"error", err,
				)
				sqsErrors[1] = err
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
	if err := serv.PropertyRepository.UpdateProperty(ctx, lite.UpdatePropertyFieldsParams{
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

func (serv *PUService) Persist(ctx context.Context, sqsBody string, userId string) error {
	sqsProcessRawImagesUrl := os.Getenv(env.SQS_PROCESS_RAW_IMAGES_QUEUE_URL)
	var propertyId = uuid.New().String()
	var requestProperty utils.RequestProperty

	if err := json.Unmarshal([]byte(sqsBody), &requestProperty); err != nil {
		serv.Log.Error("error decoding the json property", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}
	if err := json.Unmarshal([]byte(sqsBody), &requestProperty.Features); err != nil {
		serv.Log.Error("error decoding the json property features for PUT request", "err", err)
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

	if len(s3ImagesPathsProcessed) == 0 {
		// ⚠️ don't return a error
		// because I don't want lambda to retry this message
		// it's not needed to get to dlq either
		serv.Log.Error(
			"no images found in the sqs body",
			"sqs body:", requestProperty,
		)
		return nil
	}

	sqsRawImagesList := serv.SplitImagesIntoSQSImageMessages(requestProperty.ImageNames, 2, humanReadableId)
	sqsMessagesForBatch, err := serv.BatchSQSImageMessagesIntoRequestEntries(sqsRawImagesList, 10)
	if err != nil {
		serv.Log.Error(
			"error trying to create the batches of SQS Request Entries",
			"error", err,
		)
	}

	// this is going to wait
	err = serv.DispatchSQSBatches(ctx, sqsMessagesForBatch, sqsProcessRawImagesUrl)

	if err != nil {
		serv.Log.Error("error processing the sqs batches",
			"error", err,
		)
		return err
	}

	if err := serv.PropertyRepository.Add(ctx, lite.AddPropertyParams{
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

func (serv *PUService) BatchItems(originalItems []string, batchSize int) [][]string {

	var batchedSlice = make([][]string, 0)
	for i := 0; i < len(originalItems); i += batchSize {
		end := i + batchSize
		if end > len(originalItems) {
			end = len(originalItems)
		}
		batchedSlice = append(batchedSlice, originalItems[i:end])
	}

	return batchedSlice
}

// This is going to take the list of SQSRawImages messages
// Create batches of SendMessageBatchRequestEntry with MessageBody of jsoned SQSRawImages
func (serv *PUService) BatchSQSImageMessagesIntoRequestEntries(
	sqsRawImagesMessages []utils.SQSImagesMessage,
	batchSize int64,
) ([][]sqsTypes.SendMessageBatchRequestEntry, error) {
	batchesOfRequestEntrySlices := make([][]sqsTypes.SendMessageBatchRequestEntry, 0)
	for i := 0; i < len(sqsRawImagesMessages); i += int(batchSize) {
		end := i + int(batchSize)
		if end > len(sqsRawImagesMessages) {
			end = len(sqsRawImagesMessages)
		}
		batchOfRawImagesMessage := sqsRawImagesMessages[i:end]
		batchRequestEntrySlice := make([]sqsTypes.SendMessageBatchRequestEntry, 0)
		for ind, rawImageMessage := range batchOfRawImagesMessage {
			// ⚠️ handle errors
			jsonedRawImagesMessage, err := json.Marshal(rawImageMessage)
			if err != nil {
				return nil, err
			}
			batchRequestEntrySlice = append(batchRequestEntrySlice, sqsTypes.SendMessageBatchRequestEntry{
				Id:          aws.String(strconv.Itoa((ind + 1) * end)),
				MessageBody: aws.String(string(jsonedRawImagesMessage)),
			})
		}
		batchesOfRequestEntrySlices = append(batchesOfRequestEntrySlices, batchRequestEntrySlice)
	}
	return batchesOfRequestEntrySlices, nil
}

// This is going to split the images into more SQSProcessRawImages structs
// Creating more sqs messages for smaller ammout of images slices
func (serv *PUService) SplitImagesIntoSQSImageMessages(
	images []string,
	imagesBatchedSize int,
	humanReadableId string,
) []utils.SQSImagesMessage {

	batchedImages := serv.BatchItems(images, imagesBatchedSize)

	rawImagesSQSList := make([]utils.SQSImagesMessage, len(batchedImages))
	for ind, pair := range batchedImages {
		rawImagesSQSList[ind] = utils.SQSImagesMessage{
			HumanReadableId: humanReadableId,
			Images:          pair,
		}
	}
	return rawImagesSQSList
}

func (serv *PUService) DispatchSQSBatches(
	ctx context.Context,
	sqsMessagesForBatch [][]sqsTypes.SendMessageBatchRequestEntry,
	sqsUrl string,
) error {

	sqsSendBatchesErrors := make([]error, 0)
	errChan := make(chan error, len(sqsMessagesForBatch))

	var wgForBatches sync.WaitGroup
	for _, batchRequestEntry := range sqsMessagesForBatch {
		wgForBatches.Add(1)
		go func(batch []sqsTypes.SendMessageBatchRequestEntry) {
			defer wgForBatches.Done()
			_, err := serv.SQSClient.SendMessageBatch(
				ctx,
				&sqs.SendMessageBatchInput{
					QueueUrl: &sqsUrl,
					Entries:  batch,
				},
			)
			if err != nil {
				errChan <- err
			}
		}(batchRequestEntry)
	}
	wgForBatches.Wait()
	close(errChan)

	for err := range errChan {
		sqsSendBatchesErrors = append(sqsSendBatchesErrors, err)
	}

	return errors.Join(sqsSendBatchesErrors...)
}
