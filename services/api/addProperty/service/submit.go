package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"os/exec"
	"strings"

	"github.com/Serares/undertown_v3/repositories/repository"
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

func arrayToString(images []string) string {
	return strings.Join(images, ",")
}

func booleanToInt(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func (ss *Submit) ProcessPropertyImagesLocal(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var imagePaths []string
	var outErr bytes.Buffer
	err := os.Mkdir("../../ssr/homepage/assets/uploads", 0755)
	if err != nil {
		ss.Log.Error("Error creating the temp dir")
	}

	// read uploaded images
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ss.Log.Error("error trying to read the file header")
		}
		defer file.Close()
		filePath := fmt.Sprintf("%s/%s", "../../ssr/homepage/assets/uploads/", fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			ss.Log.Error("error trying to create the file path")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			ss.Log.Error("error copying the file")
		}
		webpFileName := fileHeader.Filename + ".webp"
		webpFilePath := fmt.Sprintf("%s/%s", "../../ssr/homepage/assets/uploads/", webpFileName)
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

func (ss *Submit) ProcessPropertyData(ctx context.Context, imagesPaths []string, multipartForm *multipart.Form) (string, string, error) {
	var propertyId = uuid.New().String()
	// var propertyTransaction types.TransactionType = multipartForm.Value["propertyTransaction"]
	// humanReadableId := util.HumanReadableId()
	// if err := ss.PropertyRepository.Add(ctx, lite.AddPropertyParams{
	// 	ID: propertyId,
	// 	// UserID: , TODO
	// 	Humanreadableid:                  humanReadableId,
	// 	Title:                            property.Title,
	// 	Floor:                            property.Floor,
	// 	Images:                           arrayToString(property.Images),
	// 	Thumbnail:                        property.Thumbnail,
	// 	IsFeatured:                       booleanToInt(property.IsFeatured),
	// 	EnergyClass:                      property.EnergyClass,
	// 	EnergyConsumptionPrimary:         property.EnergyConsumptionPrimary,
	// 	EnergyEmissionsIndex:             property.EnergyEmissionsIndex,
	// 	EnergyConsumptionGreen:           property.EnergyConsumptionGreen,
	// 	DestinationResidential:           booleanToInt(property.DestinationResidential),
	// 	DestinationCommercial:            booleanToInt(property.DestinationCommercial),
	// 	DestinationOffice:                booleanToInt(property.DestinationOffice),
	// 	DestinationHoliday:               booleanToInt(property.DestinationHoliday),
	// 	OtherUtilitiesTerrance:           booleanToInt(property.OtherUtilitiesTerrance),
	// 	OtherUtilitiesServiceToilet:      booleanToInt(property.OtherUtilitiesServiceToilet),
	// 	OtherUtilitiesUndergroundStorage: booleanToInt(property.OtherUtilitiesUndergroundStorage),
	// 	OtherUtilitiesStorage:            booleanToInt(property.OtherUtilitiesStorage),
	// 	PropertyTransaction:              property.PropertyTransaction.String(),
	// 	PropertyDescription:              property.PropertyDescription,
	// 	PropertyType:                     property.PropertyType,
	// 	PropertyAddress:                  property.PropertyAddress,
	// 	PropertySurface:                  int64(property.PropertySurface),
	// 	Price:                            int64(property.Price),
	// 	FurnishedNot:                     booleanToInt(property.FurnishedNot),
	// 	FurnishedPartially:               booleanToInt(property.FurnishedPartially),
	// 	FurnishedComplete:                booleanToInt(property.FurnishedComplete),
	// 	FurnishedLuxury:                  booleanToInt(property.FurnishedLuxury),
	// 	InteriorNeedsRenovation:          booleanToInt(property.InteriorNeedsRenovation),
	// 	InteriorHasRenovation:            booleanToInt(property.InteriorHasRenovation),
	// 	InteriorGoodState:                booleanToInt(property.InteriorGoodState),
	// 	HeatingTermoficare:               booleanToInt(property.HeatingTermoficare),
	// 	HeatingCentralHeating:            booleanToInt(property.HeatingCentralHeating),
	// 	HeatingBuilding:                  booleanToInt(property.HeatingBuilding),
	// 	HeatingStove:                     booleanToInt(property.HeatingStove),
	// 	HeatingRadiator:                  booleanToInt(property.HeatingRadiator),
	// 	HeatingOtherElectrical:           booleanToInt(property.HeatingOtherElectrical),
	// 	HeatingGasConvector:              booleanToInt(property.HeatingGasConvector),
	// 	HeatingInfraredPanels:            booleanToInt(property.HeatingInfraredPanels),
	// 	HeatingFloorHeating:              booleanToInt(property.HeatingFloorHeating),
	// 	CreatedAt:                        time.Now().UTC(),
	// 	UpdatedAt:                        time.Now().UTC(),
	// }); err != nil {
	// 	return "", "", fmt.Errorf("error trying to persist the order with error: %v", err)
	// }
	ss.Log.Info("the multipart form", multipartForm)
	return propertyId, "humanReadableId", nil
}
