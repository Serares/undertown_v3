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

	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
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

func (s *SubmitService) Submit(r *http.Request, authToken, humanReadableId string, isEdit bool) (lite.Property, types.PropertyFeatures, error) {
	var err error
	// TODO
	// run some validations here if needed
	// get the token from cookie
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return lite.Property{}, types.PropertyFeatures{}, fmt.Errorf("error trying to read the request body %v", err)
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
		return lite.Property{}, types.PropertyFeatures{}, fmt.Errorf("error marshaling the json structure %v", err)
	}
	_, err = textWriter.Write(jsonString)
	if err != nil {
		s.Log.Error("error writing the json string")
		return lite.Property{}, types.PropertyFeatures{}, fmt.Errorf("error writing json string to the body %v", err)
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
	// 🤔
	// If the requests fail
	// This method should return the form fields back to the view
	// to fill the values of the inputs
	writer.Close()
	if isEdit {
		url = fmt.Sprintf("%s?propertyId=%s", url, humanReadableId)
		err = s.Client.AddProperty(&newReaderBuffer, url, authToken, writer.FormDataContentType(), http.MethodPut)
	} else {
		err = s.Client.AddProperty(&newReaderBuffer, url, authToken, writer.FormDataContentType(), http.MethodPost)
	}

	if err != nil {
		// do the json unmarshal
		var property utils.RequestProperty // using the type from the addProperty service
		err := json.Unmarshal(jsonString, &property)
		if err != nil {
			return lite.Property{}, types.PropertyFeatures{}, err
		}
		// Have to unmarshal the property features into the types.PropertyFeatures struct
		// ❗but the features are unmarshaled first as a map[string]interface{}
		var features types.PropertyFeatures
		featuresAsJsonString, err := json.Marshal(property.Features)
		if err != nil {
			return lite.Property{}, types.PropertyFeatures{}, err
		}
		err = json.Unmarshal(featuresAsJsonString, &features)
		if err != nil {
			return lite.Property{}, types.PropertyFeatures{}, err
		}
		// ❗
		// populating a lite.Property struct with the types.RequestProperty fields
		// because the views.Edit() component used for rerendering the last values
		// is accepting a lite.Property prop as parameter
		// TODO think if there are other ways around this hacky sittuation
		// ❗this is the ending where the fields need to be rerendered
		return lite.Property{
			Title:               property.Title,
			Price:               property.Price,
			IsFeatured:          utils.BoolToInt(property.IsFeatured),
			PropertyType:        property.PropertyType,
			PropertyDescription: property.PropertyDescription,
			PropertyAddress:     property.PropertyAddress,
			PropertyTransaction: fmt.Sprintf("%d", property.PropertyTransaction),
			PropertySurface:     int64(property.PropertySurface),
			Features:            "", // features can be empty because it's already been unmarshaled
		}, features, nil
	}

	// this is the success ending
	return lite.Property{}, types.PropertyFeatures{}, err
}
