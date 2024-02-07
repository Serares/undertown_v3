package service

import (
	"log/slog"
	"net/http"
	"os"

	adminUtils "github.com/Serares/ssr/admin/utils"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
)

type SubmitService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewSubmitService(log *slog.Logger, client ISSRAdminClient) *SubmitService {
	return &SubmitService{
		Log:    log.WithGroup("Submit Service"),
		Client: client,
	}
}

type PropertyFormField struct {
	Title string
}

func (s *SubmitService) Submit(r *http.Request, authToken, humanReadableId string) (lite.Property, utils.PropertyFeatures, error) {
	var err error
	// TODO
	// run some validations here if needed
	// get the token from cookie
	url := os.Getenv("SUBMIT_PROPERTY_URL")
	// 🤔
	// If the requests fail
	// This method should return the form fields back to the view
	// to fill the values of the inputs
	bufferedBody, contentType, jsonString, err := adminUtils.ParseMultipart(r)
	if err != nil {
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	err = s.Client.AddProperty(bufferedBody, url, authToken, contentType, http.MethodPost)

	if err != nil {
		s.Log.Error("error trying to send the request", "error", err)
		property, features, unmarshallingErrors := adminUtils.UnmarshalProperty(jsonString)
		// ❗
		// populating a lite.Property struct with the types.RequestProperty fields
		// because the views.Edit() component used for rerendering the last values
		// is accepting a lite.Property prop as parameter
		// TODO think if there are other ways around this hacky sittuation
		// ❗this is the ending where the fields need to be rerendered
		if unmarshallingErrors != nil {
			s.Log.Error("error trying to unmarshal the data for error response", "error", unmarshallingErrors)
			// This case the property and features will be empty struct values
			return property, features, err
		}
		return property, features, err
	}

	// this is the success ending
	return lite.Property{}, utils.PropertyFeatures{}, nil
}
