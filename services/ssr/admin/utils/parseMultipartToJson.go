package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ❗TODO
//
//	Split the S3 upload logic into a separate function
func ParseMultipartToJson(r *http.Request) ([]byte, string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, "", err
	}

	var jsonStructure map[string]interface{} = make(map[string]interface{})
	for key, values := range r.PostForm {
		if len(values) > 1 {
			// this is for all the fields that have one key with multiple values
			jsonStructure[key] = values
		} else {
			var err error
			var number int64
			// Value[0] it's creating arrays of input values
			// in case there are two inputs with the same name
			// it will create an array with the same key and more values
			// check if int
			if number, err = strconv.ParseInt(values[0], 10, 64); err == nil {
				jsonStructure[key] = number
			}
			// check if checkbox or string
			if err != nil {
				if values[0] == "on" {
					jsonStructure[key] = true
				} else {
					jsonStructure[key] = values[0]
				}
			}
		}
	}

	var imagesNamesList = make([]string, 0)
	cfg, err := config.LoadDefaultConfig(r.Context())
	if err != nil {
		return nil, "", err
	}

	s3Client := s3.NewFromConfig(cfg)
	// Get the files from the multipart form and persist them to S3
	// Add the file names to the json string

	// ❗TODO
	// find a different way of getting the property_transaction from the form
	propertyTransactionFormValue := r.PostForm.Get(constants.TransactionTypeFormInputKey)
	propertyTransactionToInt, err := strconv.Atoi(propertyTransactionFormValue)
	if err != nil {
		return nil, "", err
	}
	transactionType := utils.TransactionType(propertyTransactionToInt)

	humanReadableId := utils.HumanReadableId(
		transactionType,
	)

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			imageName := utils.ReplaceWhiteSpaceWithUnderscore(
				fmt.Sprintf(
					"%s/%s_%s",
					humanReadableId,
					utils.GenerateStringTimestamp(),
					fileHeader.Filename,
				),
			)
			bucketKey := "/" + imageName
			file, err := fileHeader.Open()

			if err != nil {
				return nil, "", fmt.Errorf("error reading the file from the form %v", err)
			}

			defer file.Close()
			imagesNamesList = append(imagesNamesList, imageName)
			// ❗
			// TODO run this in goroutines after checking it works
			err = UploadFilesToS3(
				r.Context(),
				*fileHeader,
				s3Client,
				bucketKey,
			)
			if err != nil {
				return nil, "", fmt.Errorf("error uploading the file '%s' to s3 %v", bucketKey, err)
			}
		}
	}

	jsonStructure[constants.ImagesFormKey] = imagesNamesList

	// json marshal
	jsonString, err := json.Marshal(jsonStructure)
	if err != nil {
		return nil, "", fmt.Errorf("error marshaling the json structure %v", err)
	}

	return jsonString, humanReadableId, err
}
