package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
)

type PropertyService struct {
	Log    *slog.Logger
	Client ISSRClient
}

func NewPropertyService(log *slog.Logger, client ISSRClient) *PropertyService {
	return &PropertyService{
		Log:    log,
		Client: client,
	}
}

func (s *PropertyService) Get(humanReadableId string) (types.ProcessedSingleProperty, error) {
	getPropertyUrl := os.Getenv("GET_PROPERTY_URL")
	url, err := utils.AddParamToUrl(getPropertyUrl, constants.HumanReadableIdQueryKey, humanReadableId)
	if err != nil {
		return types.ProcessedSingleProperty{}, fmt.Errorf("error trying to create the get property url %v", err)
	}
	liteProperty, err := s.Client.GetProperty(url)
	if err != nil {
		return types.ProcessedSingleProperty{}, fmt.Errorf("error trying to get the property from the backend %v", err)
	}

	// process the property
	var processedProperty types.ProcessedSingleProperty

	// unmarshal the features
	err = json.Unmarshal([]byte(liteProperty.Features), &processedProperty.Features)
	if err != nil {
		return types.ProcessedSingleProperty{}, fmt.Errorf("error trying to unmarshal the property features into the struct %v", err)
	}

	processedProperty.Title = liteProperty.Title
	processedProperty.Address = liteProperty.PropertyAddress
	processedProperty.Description = liteProperty.PropertyDescription
	processedProperty.DisplayPrice = utils.CreateDisplayPrice(liteProperty.Price)
	processedProperty.ImagePaths = utils.CreateImagePathList("/images/", strings.Split(liteProperty.Images, ";"))
	processedProperty.Surface = fmt.Sprintf("%d", liteProperty.PropertySurface)

	return processedProperty, nil
}
