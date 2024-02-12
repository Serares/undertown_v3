package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"log/slog"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	repositoryTypes "github.com/Serares/undertown_v3/repositories/repository/types"
	"github.com/Serares/undertown_v3/services/api/addProperty/types"
	"github.com/Serares/undertown_v3/services/api/addProperty/util"
	"github.com/Serares/undertown_v3/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

var LocalAssetsRelativePath = "../../ssr/assets/uploads"

// TODO handle the property creation
// handle the checking if the user id exists before persisting the property

// TODO make sure all fields are mapped from repository model
// TODO move this struct to the handler
type Submit struct {
	Log                *slog.Logger
	PropertyRepository *repository.Properties
}

func NewSubmitService(log *slog.Logger, pr *repository.Properties) Submit {
	return Submit{
		Log:                log.WithGroup("Submit Service"),
		PropertyRepository: pr,
	}
}

func (ss *Submit) encodeToWebP(file multipart.File) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	img, err := jpeg.Decode(file)
	if err != nil {
		ss.Log.Error("error deconding the file",
			"error", err,
		)
	}

	err = webp.Encode(&buf, img, &webp.Options{Lossless: true})
	if err != nil {
		ss.Log.Error("error encoding the file",
			"error", err,
		)
	}

	return &buf, err
}

func (ss *Submit) ProcessImagesLocal(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var imageNames []string
	err := os.Mkdir(LocalAssetsRelativePath, 0755)
	if err != nil {
		ss.Log.Error("Error creating the temp dir", "err", err)
	}

	// read uploaded images
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ss.Log.Error("error trying to read the file header")
		}
		defer file.Close()
		fileName := utils.ReplaceWhiteSpaceWithUnderscore(fileHeader.Filename)

		buf, err := ss.encodeToWebP(file)
		if err != nil {
			ss.Log.Error("error trying to encode to webp",
				"error", err,
			)
		}
		webpFileName := LocalAssetsRelativePath + "/" + fileName + ".webp"
		webpFile, err := os.Create(webpFileName)
		if err != nil {
			ss.Log.Error("error creating webFile the file",
				"error", err,
			)
		}
		defer webpFile.Close()

		_, err = webpFile.Write(buf.Bytes())
		if err != nil {
			ss.Log.Error("error writing webFile the file",
				"error", err,
			)
		}
		imageNames = append(imageNames, webpFileName)
		// remove files after beeing uploaded to s3 to not overflow to storage
	}
	return imageNames, nil
}

// Uploads the files to S3
func (ss *Submit) ProcessImagesS3(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var imageNames []string
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		ss.Log.Error("error loading the default config for lambda", "error", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	bucketName := os.Getenv(types.ASSETS_BUCKET_ENV)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ss.Log.Error("error trying to read the file header")
		}
		defer file.Close()
		fileName := utils.ReplaceWhiteSpaceWithUnderscore(fileHeader.Filename)
		buf, err := ss.encodeToWebP(file)
		if err != nil {
			ss.Log.Error("error trying to encode to webp",
				"error", err,
			)
		}
		webpFileName := fileName + ".webp"

		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("uploads/" + webpFileName),
			Body:   buf,
		})
		if err != nil {
			ss.Log.Error("error uploading the file to s3", "error", err)
		}
		imageNames = append(imageNames, webpFileName)
		// remove files after beeing uploaded to s3 to not overflow to storage
	}
	return imageNames, nil
}

func (ss *Submit) parsePropertyFeaturesToJson(features types.RequestFeatures) (string, error) {
	featuresJson, err := json.Marshal(features)
	if err != nil {
		ss.Log.Error("error marshalling the features to json string")
		return "", err
	}
	return string(featuresJson), nil
}

func (ss *Submit) ProcessPropertyUpdateData(ctx context.Context, imagesPaths, deleteImages []string, multipartForm *multipart.Form, humanReadableId string) error {
	var requestProperty utils.RequestProperty
	jsonProperty, ok := multipartForm.Value["property"]
	if !ok {
		ss.Log.Error("json property not provided for PUT request")
		return fmt.Errorf("json property not provided")
	}

	if err := json.Unmarshal([]byte(jsonProperty[0]), &requestProperty); err != nil {
		ss.Log.Error("error decoding the json property for PUT request", "err", err)
		return fmt.Errorf("error on json unmarshal")
	}
	if err := json.Unmarshal([]byte(jsonProperty[0]), &requestProperty.Features); err != nil {
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
	if len(imagesPaths) > 0 {
		finalImages = append(finalImages, imagesPaths...)
	}
	var filteredImages = make([]string, 0)
	removeMap := make(map[string]bool, 0)
	for _, img := range deleteImages {
		removeMap[img] = true
	}

	for _, img := range finalImages {
		if !removeMap[img] {
			filteredImages = append(filteredImages, img)
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
		PropertyTransaction: repositoryTypes.TransactionType(requestProperty.PropertyTransaction).String(),
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

func (ss *Submit) ProcessPropertyData(ctx context.Context, imagesPaths []string, multipartForm *multipart.Form, userId string) (string, string, error) {
	var propertyId = uuid.New().String()
	var requestProperty utils.RequestProperty
	thumbnail := ""
	jsonProperty, ok := multipartForm.Value["property"]
	if !ok {
		ss.Log.Error("json property not provided")
		return "", "", fmt.Errorf("json property not provided")
	}

	if err := json.Unmarshal([]byte(jsonProperty[0]), &requestProperty); err != nil {
		ss.Log.Error("error decoding the json property", "err", err)
		return "", "", fmt.Errorf("error on json unmarshal")
	}
	if err := json.Unmarshal([]byte(jsonProperty[0]), &requestProperty.Features); err != nil {
		ss.Log.Error("error decoding the json property features for PUT request", "err", err)
		return "", "", fmt.Errorf("error on json unmarshal")
	}

	// SEE lite.Property
	// features is stored as a json string in the db
	// propertyFeatures := types.RequestFeatures{
	// 	Floor:                            requestProperty.Floor,
	// 	EnergyClass:                      requestProperty.EnergyClass,
	// 	EnergyConsumptionPrimary:         requestProperty.EnergyConsumptionPrimary,
	// 	EnergyEmissionsIndex:             requestProperty.EnergyEmissionsIndex,
	// 	EnergyConsumptionGreen:           requestProperty.EnergyConsumptionGreen,
	// 	DestinationResidential:           requestProperty.DestinationResidential,
	// 	DestinationCommercial:            requestProperty.DestinationCommercial,
	// 	DestinationOffice:                requestProperty.DestinationOffice,
	// 	DestinationHoliday:               requestProperty.DestinationHoliday,
	// 	OtherUtilitiesTerrance:           requestProperty.OtherUtilitiesTerrance,
	// 	OtherUtilitiesServiceToilet:      requestProperty.OtherUtilitiesServiceToilet,
	// 	OtherUtilitiesUndergroundStorage: requestProperty.OtherUtilitiesUndergroundStorage,
	// 	OtherUtilitiesStorage:            requestProperty.OtherUtilitiesStorage,
	// 	FurnishedNot:                     requestProperty.FurnishedNot,
	// 	FurnishedPartially:               requestProperty.FurnishedPartially,
	// 	FurnishedComplete:                requestProperty.FurnishedComplete,
	// 	FurnishedLuxury:                  requestProperty.FurnishedLuxury,
	// 	InteriorNeedsRenovation:          requestProperty.InteriorNeedsRenovation,
	// 	InteriorHasRenovation:            requestProperty.InteriorHasRenovation,
	// 	InteriorGoodState:                requestProperty.InteriorGoodState,
	// 	HeatingTermoficare:               requestProperty.HeatingTermoficare,
	// 	HeatingCentralHeating:            requestProperty.HeatingCentralHeating,
	// 	HeatingBuilding:                  requestProperty.HeatingBuilding,
	// 	HeatingStove:                     requestProperty.HeatingStove,
	// 	HeatingRadiator:                  requestProperty.HeatingRadiator,
	// 	HeatingOtherElectrical:           requestProperty.HeatingOtherElectrical,
	// 	HeatingGasConvector:              requestProperty.HeatingGasConvector,
	// 	HeatingInfraredPanels:            requestProperty.HeatingInfraredPanels,
	// 	HeatingFloorHeating:              requestProperty.HeatingFloorHeating,
	// }

	// features, err := ss.parsePropertyFeaturesToJson(propertyFeatures)
	features, err := json.Marshal(requestProperty.Features)
	if err != nil {
		return "", "", err
	}

	if len(imagesPaths) > 0 {
		thumbnail = imagesPaths[0]
	}

	humanReadableId := util.HumanReadableId(repositoryTypes.TransactionType(requestProperty.PropertyTransaction))
	if err := ss.PropertyRepository.Add(ctx, lite.AddPropertyParams{
		ID:                  propertyId,
		UserID:              userId,
		Humanreadableid:     humanReadableId,
		Title:               requestProperty.Title,
		Images:              strings.Join(imagesPaths, ";"),
		Thumbnail:           thumbnail,
		IsFeatured:          utils.BoolToInt(requestProperty.IsFeatured),
		PropertyTransaction: repositoryTypes.TransactionType(requestProperty.PropertyTransaction).String(),
		PropertyDescription: requestProperty.PropertyDescription,
		PropertyType:        requestProperty.PropertyType,
		PropertyAddress:     requestProperty.PropertyAddress,
		PropertySurface:     int64(requestProperty.PropertySurface),
		Price:               int64(requestProperty.Price),
		Features:            string(features),
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}); err != nil {
		return "", "", fmt.Errorf("error trying to persist the order with error: %v", err)
	}
	return propertyId, humanReadableId, nil
}

func (ss *Submit) DeleteImagesLocal(imagesNames []string) error {
	for _, imageName := range imagesNames {
		imageRelativePath := fmt.Sprintf("%s/%s", LocalAssetsRelativePath, imageName)
		if err := os.Remove(imageRelativePath); err != nil {
			ss.Log.Error("error trying to remove the file", "filePath", imageRelativePath, "error", err)
			// return fmt.Errorf("error trying to remove image %v", err) ❗this is not really needed to return because the error will generally be when files are non existent
		}
	}
	return nil
}

func (ss *Submit) DeleteImagesS3(ctx context.Context, imagesList []string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	bucketName := os.Getenv(types.ASSETS_BUCKET_ENV)
	if err != nil {
		ss.Log.Error("error trying to initialize the default lambda config", "error", err)
	}

	client := s3.NewFromConfig(cfg)
	var objects []s3Types.ObjectIdentifier
	for _, imageName := range imagesList {
		s3Key := fmt.Sprintf("uploads/%s", imageName)
		objects = append(objects, s3Types.ObjectIdentifier{
			Key: &s3Key,
		})
	}

	_, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: &bucketName,
		Delete: &s3Types.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		ss.Log.Error("error trying to delete objects in s3", "error", err)
		return err
	}

	return nil
}

func (s *Submit) ProcessDeleteAndPersistImages(ctx context.Context, files []*multipart.FileHeader, imagesToRemove []string) ([]string, error) {
	isLocal := os.Getenv("IS_LOCAL")
	var imageProcessErrors []error = make([]error, 2)
	var imageProcessWg sync.WaitGroup
	var imageNames []string
	// ❗
	// There might be issues if the processing takes longer than lambda execution context
	// TODO, use channels and a select block to implement a timeout
	/**
	done := make(chan struct{})
	// Goroutine to close the done channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(done)
	}()

	// Implement timeout logic
	select {
	case <-done:
		// All goroutines have finished
		fmt.Println("All goroutines finished successfully")
	case <-time.After(3 * time.Second):
		// Timeout occurred
		fmt.Println("Timeout occurred waiting for goroutines to finish")
	}
	**/
	imageProcessWg.Add(2)
	go func() {
		defer imageProcessWg.Done()
		if len(files) > 0 {
			if isLocal == "true" {
				imageNames, imageProcessErrors[0] = s.ProcessImagesLocal(ctx, files)
			} else {
				imageNames, imageProcessErrors[0] = s.ProcessImagesS3(ctx, files)
			}
		}
	}()

	go func() {
		defer imageProcessWg.Done()
		if len(imagesToRemove) > 0 {
			if isLocal == "true" {
				imageProcessErrors[1] = s.DeleteImagesLocal(imagesToRemove)
			} else {
				imageProcessErrors[1] = s.DeleteImagesS3(ctx, imagesToRemove)
			}
		}
	}()

	imageProcessWg.Wait()

	return imageNames, errors.Join(imageProcessErrors...)
}
