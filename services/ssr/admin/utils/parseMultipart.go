package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// ❔
// this is going to return the buffered body
// the content type of the multipart/form request
// the json string used to rerender the data in case something fails
// and error
func ParseMultipart(r *http.Request) (*bytes.Buffer, string, []byte, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, "", nil, err
	}

	var newReaderBuffer bytes.Buffer
	writer := multipart.NewWriter(&newReaderBuffer)
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

	textWriter, _ := writer.CreateFormField("property")
	// json marshal
	jsonString, err := json.Marshal(jsonStructure)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error marshaling the json structure %v", err)
	}
	_, err = textWriter.Write(jsonString)
	if err != nil {
		return nil, "", jsonString, fmt.Errorf("error writing json string to the body %v", err)
	}
	// get the files from the multipar form
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, "", nil, fmt.Errorf("error reading the file from the form %v", err)
			}
			defer file.Close()

			fw, err := writer.CreateFormFile("images", fileHeader.Filename)
			if err != nil {
				return nil, "", nil, fmt.Errorf("error creating file writer %v", err)
			}
			if _, err = io.Copy(fw, file); err != nil {
				return nil, "", nil, fmt.Errorf("error writing the form file to the request multipart form file %v", err)
			}
		}
	}
	contentType := writer.FormDataContentType()
	writer.Close()
	return &newReaderBuffer, contentType, jsonString, err
}
