package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"os/exec"
	"strings"
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
	"github.com/google/uuid"
)

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
		Log:                log,
		PropertyRepository: pr,
	}
}

func (ss *Submit) ProcessPropertyImagesLocal(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var localDirPath = "../../ssr/assets/uploads"
	var imagePaths []string
	var outErr bytes.Buffer
	err := os.Mkdir(localDirPath, 0755)
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
		filePath := fmt.Sprintf("%s/%s", localDirPath, fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			ss.Log.Error("error trying to create the file path")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			ss.Log.Error("error copying the file")
		}
		webpFileName := fileHeader.Filename + ".webp"
		webpFilePath := fmt.Sprintf("%s/%s", localDirPath, webpFileName)
		// convert to webp
		cmd := exec.Command("cwebp", "-q", "75", filePath, "-o", webpFilePath)
		cmd.Stderr = &outErr
		if err := cmd.Run(); err != nil {
			ss.Log.Error("error converting the file", "cmd error:", outErr.String())
		}
		imagePaths = append(imagePaths, "/uploads/"+webpFileName)
		// remove files after beeing uploaded to s3 to not overflow to storage
		defer os.Remove(filePath)
	}
	return imagePaths, nil
}

// Uploads the files to S3
func (ss *Submit) ProcessPropertyImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var imagePaths []string
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		ss.Log.Error("error loading the default config for lambda")
	}
	s3Client := s3.NewFromConfig(cfg)
	var outErr bytes.Buffer
	bucketName := os.Getenv("ASSETS_BUCKET_NAME")
	webpDir, err := os.MkdirTemp("", "webp")
	if err != nil {
		ss.Log.Error("error creating webp dir")
	}
	formImagesTempDir, err := os.MkdirTemp("", "images")
	if err != nil {
		ss.Log.Error("Error creating the temp dir")
	}

	defer os.RemoveAll(formImagesTempDir)
	defer os.RemoveAll(webpDir)
	// read uploaded images
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ss.Log.Error("error trying to read the file header")
		}
		defer file.Close()
		filePath := fmt.Sprintf("%s/%s", formImagesTempDir, fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			ss.Log.Error("error trying to create the file path")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			ss.Log.Error("error copying the file")
		}
		webpFileName := fileHeader.Filename + ".webp"
		webpFilePath := fmt.Sprintf("%s/%s", webpDir, webpFileName)
		// convert to webp
		cmd := exec.Command("cwebp", "-q", "75", filePath, "-o", webpFilePath)
		cmd.Stderr = &outErr
		if err := cmd.Run(); err != nil {
			ss.Log.Error("error converting the file", "cmd error:", outErr.String())
		}

		// read the webp generated file
		webpFile, err := os.Open(webpFilePath)
		if err != nil {
			ss.Log.Error("error opening the webp file path:", "path", webpFilePath)
		}
		defer webpFile.Close()
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("uploads/" + webpFileName),
			Body:   webpFile,
		})
		if err != nil {
			ss.Log.Error("error uploading the file to s3", "error", err)
		}
		imagePaths = append(imagePaths, "/uploads/"+webpFileName)
		// remove files after beeing uploaded to s3 to not overflow to storage
		defer os.Remove(webpFilePath)
		defer os.Remove(filePath)
	}
	return imagePaths, nil
}

func (ss *Submit) parsePropertyFeaturesToJson(features types.RequestFeatures) (string, error) {
	featuresJson, err := json.Marshal(features)
	if err != nil {
		ss.Log.Error("error marshalling the features to json string")
		return "", err
	}
	return string(featuresJson), nil
}

func (ss *Submit) ProcessPropertyUpdateData(ctx context.Context, imagesPaths []string, multipartForm *multipart.Form, humanReadableId string) error {
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

	// TODO handle the case where no new images are uploaded
	// OR when new images are uploaded and the old ones are not deleted
	if err := ss.PropertyRepository.UpdateProperty(ctx, lite.UpdatePropertyFieldsParams{
		Humanreadableid:     humanReadableId,
		Title:               requestProperty.Title,
		Images:              strings.Join(imagesPaths, ";"),
		Thumbnail:           imagesPaths[0],
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

	humanReadableId := util.HumanReadableId(repositoryTypes.TransactionType(requestProperty.PropertyTransaction))
	if err := ss.PropertyRepository.Add(ctx, lite.AddPropertyParams{
		ID:                  propertyId,
		UserID:              userId,
		Humanreadableid:     humanReadableId,
		Title:               requestProperty.Title,
		Images:              strings.Join(imagesPaths, ";"),
		Thumbnail:           imagesPaths[0],
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
