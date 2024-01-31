package service

import (
	"fmt"
	"log/slog"
	"os"

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
		tempLp := types.ListingProperty{
			Title:              property.Title,
			Price:              property.Price,
			Address:            property.PropertyAddress,
			TransactionType:    property.PropertyTransaction,
			DisplayPrice:       utils.CreateDisplayPrice(property.Price),
			Thumbnail:          property.Thumbnail,
			EditPropertyPath:   fmt.Sprintf("/edit/%s", utils.CreatePropertyPath(property.Title, property.Humanreadableid)),
			DeletePropertyPath: fmt.Sprintf("/delete/%s", utils.CreatePropertyPath(property.Title, property.Humanreadableid)),
			Surface:            property.PropertySurface,
			ImagesNumber:       2, // TODO the images are a string separated by ;
			CreatedTime:        utils.CreateDisplayCreatedAt(property.CreatedAt),
		}
		listingProperties = append(listingProperties, tempLp)
	}

	return listingProperties, nil
}
