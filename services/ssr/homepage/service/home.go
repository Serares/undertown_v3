package service

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/undertown_v3/utils"
)

type HomeService struct {
	Log    *slog.Logger
	client ISSRClient
}

// TODO does the return have to be a pointer?
func NewHomeService(log *slog.Logger, client ISSRClient) *HomeService {
	return &HomeService{
		Log:    log,
		client: client,
	}
}

func (hs *HomeService) Get() ([]types.ProcessedFeaturedProperty, error) {
	getPropertiesUrl := os.Getenv("GET_PROPERTIES_URL")
	var processedFeatProperties []types.ProcessedFeaturedProperty

	properties, err := hs.client.ListFeaturedProperties(strings.Join([]string{getPropertiesUrl, "featured=true"}, "?"))
	if err != nil {
		return []types.ProcessedFeaturedProperty{}, err
	}

	for _, featProp := range properties {
		propertyThumbnailPath := utils.CreateImagePath("/images/", featProp.Thumbnail)

		propertyPath, err := utils.CreateSinglePropertyPath(featProp.PropertyTransaction, featProp.Title, featProp.Humanreadableid)
		if err != nil {
			hs.Log.Error("error creating the property path", "error", err)
			return []types.ProcessedFeaturedProperty{}, fmt.Errorf("error trying to generate the property path %v", err)
		}
		processedFeatProperties = append(processedFeatProperties, types.ProcessedFeaturedProperty{
			Title:           featProp.Title,
			TransactionType: featProp.PropertyTransaction,
			PropertyType:    featProp.PropertyType, // TODO does it need to be processed?
			Price:           featProp.Price,
			DisplayPrice:    utils.CreateDisplayPrice(featProp.Price),
			PropertyPathUrl: propertyPath,
			CreatedTime:     utils.CreateDisplayCreatedAt(featProp.CreatedAt),
			ThumbnailPath:   propertyThumbnailPath,
		})
	}
	return processedFeatProperties, nil
}
