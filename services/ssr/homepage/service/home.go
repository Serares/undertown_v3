package service

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/undertown_v3/utils"
)

type ProcessedFeaturedProperty struct {
	Title           string
	TransactionType string
	Price           int64
	DisplayPrice    string
	Thumbnail       string
	PropertyPathUrl string
	CreatedTime     string
}

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

func (hs *HomeService) ListProperties() ([]ProcessedFeaturedProperty, error) {
	getPropertiesUrl := os.Getenv("GET_PROPERTIES_URL")
	var processedFeatProperties []ProcessedFeaturedProperty

	properties, err := hs.client.ListFeaturedProperties(strings.Join([]string{getPropertiesUrl, "featured=true"}, "?"))
	if err != nil {
		return []ProcessedFeaturedProperty{}, err
	}

	for _, featProp := range properties {
		processedFeatProperties = append(processedFeatProperties, ProcessedFeaturedProperty{
			Title:           featProp.Title,
			TransactionType: featProp.PropertyTransaction,
			Price:           featProp.Price,
			DisplayPrice:    utils.CreateDisplayPrice(featProp.Price),
			PropertyPathUrl: utils.CreatePropertyPath(featProp.Title, featProp.Humanreadableid),
			CreatedTime:     utils.CreateDisplayCreatedAt(featProp.CreatedAt),
			Thumbnail:       featProp.Thumbnail,
		})
	}
	return processedFeatProperties, nil
}
