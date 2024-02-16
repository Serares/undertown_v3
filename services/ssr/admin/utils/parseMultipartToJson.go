package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	rootUtils "github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ParseMultipartFieldsToJson(
	r *http.Request,
	humanReadableId string,
	s3Client *s3.Client,
) ([]byte, error) {
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

	// Get the files from the multipart form and persist them to S3
	// Add the file names to the json string

	var wg sync.WaitGroup
	// ❗TODO
	// find a different way of getting the property_transaction from the form
	var fileParsingErrors = make([]error, 0)
	var s3UploadingErrors = make([]error, 0)
	errChan := make(chan error, len(r.MultipartForm.File))

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			bucketKey := rootUtils.ReplaceWhiteSpaceWithUnderscore(
				fmt.Sprintf(
					"%s/%s_%s",
					humanReadableId,
					rootUtils.GenerateStringTimestamp(),
					fileHeader.Filename,
				),
			)
			file, err := fileHeader.Open()
			if err != nil {
				fileParsingErrors = append(fileParsingErrors, err)
				continue
			}
			defer file.Close()
			wg.Add(1)
			go UploadFilesToS3(
				r.Context(),
				file,
				s3Client,
				bucketKey,
				&wg,
				errChan,
			)
			imagesNamesList = append(imagesNamesList, bucketKey)
		}
	}

	wg.Wait()
	close(errChan)
	for s3Error := range errChan {
		if s3Error != nil {
			s3UploadingErrors = append(s3UploadingErrors, s3Error)
		}
	}

	if len(s3UploadingErrors) > 0 {
		return nil, errors.Join(s3UploadingErrors...)
	}
	if len(fileParsingErrors) > 0 {
		return nil, errors.Join(fileParsingErrors...)
	}
	// 💩 looks like a shitty pattern
	// but this is the place that's handling S3 uploading and images names creation
	if deletedImagesSlice, ok := jsonStructure[constants.DeleteImagesFormKey].([]string); ok {
		// append the human readable ID to the images that will be deleted
		for i, delImgName := range deletedImagesSlice {
			deletedImagesSlice[i] = humanReadableId + "/" + delImgName
		}
	}
	// 💩 this is poopoo
	// because if there is only one deleted_image the multipart/form will not send the values as a slice
	if deletedImage, ok := jsonStructure[constants.DeleteImagesFormKey].(string); ok {
		deletdImageHrId := humanReadableId + "/" + deletedImage
		jsonStructure[constants.DeleteImagesFormKey] = []string{deletdImageHrId}
	}

	jsonStructure[constants.ImagesFormKey] = imagesNamesList

	// json marshal
	jsonString, err := json.Marshal(jsonStructure)
	if err != nil {
		return nil, fmt.Errorf("error marshaling the json structure %v", err)
	}

	return jsonString, err
}
