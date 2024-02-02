package service

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/undertown_v3/utils"
)

type ListingsService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewListingService(log *slog.Logger, client ISSRAdminClient) *ListingsService {
	return &ListingsService{
		Log:    log,
		Client: client,
	}
}

func (ls *ListingsService) List(authToken string) ([]types.ListingProperty, error) {
	getPropertiesUrl := os.Getenv("GET_PROPERTIES_URL")
	// return only the necessary fields to the views
	properties, err := ls.Client.List(getPropertiesUrl, authToken)
	if err != nil {
		return []types.ListingProperty{}, err
	}
	listingProperties := make([]types.ListingProperty, 0)
	for _, property := range properties {
		editUrl := fmt.Sprintf("%s/%s", types.EditPath, utils.UrlEncodePropertyTitle(property.Title))
		// add property human readable id as query string
		editUrl, err = utils.AddParamToUrl(editUrl, utils.HumanReadableIdQueryKey, property.Humanreadableid)
		if err != nil {
			ls.Log.Error("error creating the edit url for proeprty:", "hrID", property.Humanreadableid)
			editUrl = fmt.Sprintf("%s/brokenurl", types.EditPath)
		}
		imagesLen := len(strings.Split(property.Images, ";"))
		tempLp := types.ListingProperty{
			Title:            property.Title,
			Price:            property.Price,
			Address:          property.PropertyAddress,
			TransactionType:  property.PropertyTransaction,
			DisplayPrice:     utils.CreateDisplayPrice(property.Price),
			Thumbnail:        property.Thumbnail,
			EditPropertyPath: editUrl,
			Surface:          property.PropertySurface,
			ImagesNumber:     int64(imagesLen),
			CreatedTime:      utils.CreateDisplayCreatedAt(property.CreatedAt),
		}
		listingProperties = append(listingProperties, tempLp)
	}

	return listingProperties, nil
}
