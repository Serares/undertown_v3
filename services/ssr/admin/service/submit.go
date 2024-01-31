package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type SubmitService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewSubmitService(log *slog.Logger, client ISSRAdminClient) *SubmitService {
	return &SubmitService{
		Log:    log,
		Client: client,
	}
}

type PropertyFormField struct {
	Title string
}

func (s *SubmitService) Submit(r *http.Request, authToken, humanReadableId string, isEdit bool) error {
	// TODO
	// run some validations here if needed
	// get the token from cookie
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return fmt.Errorf("error trying to read the request body %v", err)
	}
	var newReaderBuffer bytes.Buffer
	writer := multipart.NewWriter(&newReaderBuffer)
	var jsonStructure map[string]interface{} = make(map[string]interface{})

	for key, values := range r.PostForm {
		for _, value := range values {
			var err error
			var number int64
			// Value[0] it's creating arrays of input values
			// in case there are two inputs with the same name
			// it will create an array with the same key and more values
			// check if int
			if number, err = strconv.ParseInt(value, 10, 64); err == nil {
				jsonStructure[key] = number
			}
			// check if checkbox or string
			if err != nil {
				if value == "on" {
					jsonStructure[key] = true
				} else {
					jsonStructure[key] = value
				}
			}
		}
	}

	textWriter, _ := writer.CreateFormField("property")
	// json marshal
	jsonString, err := json.Marshal(jsonStructure)
	if err != nil {
		s.Log.Error("error writing the json string")
		return fmt.Errorf("error marshaling the json structure %v", err)
	}
	_, err = textWriter.Write(jsonString)
	if err != nil {
		s.Log.Error("error writing the json string")
		return fmt.Errorf("error writing json string to the body %v", err)
	}
	// get the files from the multipar form
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				s.Log.Error("error reading the file from the form", "err", err)
			}
			defer file.Close()

			fw, err := writer.CreateFormFile("images", fileHeader.Filename)
			if err != nil {
				s.Log.Error("error creating a file writer for form", "err", err)
			}
			if _, err = io.Copy(fw, file); err != nil {
				s.Log.Error("error writing the form file to the request multipart form file", "err", err)
			}
		}
	}
	url := os.Getenv("SUBMIT_PROPERTY_URL")
	// client := &http.Client{}
	// request, err := http.NewRequest(http.MethodPost, submitUrl, &newReaderBuffer)
	// if err != nil {
	// 	s.Log.Error("error on creating the request", "err", err)
	// }
	// request.Header.Set("Content-Type", writer.FormDataContentType())
	// // request.Header.Set("Authentication", authToken)
	// writer.Close()
	// response, err := client.Do(request)
	// if err != nil {
	// 	s.Log.Error("error on sending the request", "err", err)
	// }
	// defer response.Body.Close()
	// byteResponse, err := io.ReadAll(response.Body)
	// s.Log.Info("response body", "response", string(byteResponse))
	writer.Close()
	if isEdit {
		url = fmt.Sprintf("%s?propertyId=%s", url, humanReadableId)
		return s.Client.AddProperty(&newReaderBuffer, url, authToken, writer.FormDataContentType(), http.MethodPut)
	} else {
		return s.Client.AddProperty(&newReaderBuffer, url, authToken, writer.FormDataContentType(), http.MethodPost)
	}
}
