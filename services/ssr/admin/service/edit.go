package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type EditService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewEditService(log *slog.Logger, client ISSRAdminClient) *EditService {
	return &EditService{
		Log:    log,
		Client: client,
	}
}

// TODO right now the Submit() method from the SubmitService is used to edit the property
// func (es *EditService) Post(body, humanReadableId, authToken string) error {
// }

func (es *EditService) Get(humanReadableId, authToken string) (lite.Property, []string, types.PropertyFeatures, error) {
	getPropertyUrl := os.Getenv("GET_PROPERTY_URL")
	// have to add the human readable id to the url
	// ❗TODO
	// this might need some validations
	urlWithQueryString := fmt.Sprintf("%s?propertyId=%s", getPropertyUrl, humanReadableId)

	// process the images returned from the db
	// the images are a string separated by ;
	// return to the views the images paths as a slice of strings

	property, err := es.Client.GetProperty(urlWithQueryString, authToken)
	if err != nil {
		return lite.Property{}, nil, types.PropertyFeatures{}, err
	}

	images := strings.Split(property.Images, ";")
	// have to construct the path with /assets/ as a prefix
	imagesWithPrefix := make([]string, len(images))
	for idx, img := range images {
		imagesWithPrefix[idx] = types.ImagesPrefix + img
	}

	// have to decode the property features into the types.PropertyFeatures struct to be able to fill
	// up input values
	var propertyFeatures types.PropertyFeatures

	err = json.Unmarshal([]byte(property.Features), &propertyFeatures)
	if err != nil {
		return lite.Property{}, nil, types.PropertyFeatures{}, err
	}

	return property, imagesWithPrefix, propertyFeatures, nil
}
