package service

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/Serares/ssr/homepage/types"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	repositoryTypes "github.com/Serares/undertown_v3/repositories/repository/types"
	rootUtils "github.com/Serares/undertown_v3/utils"
)

const (
	Sell = "vanzari"
	Rent = "chirii"
)

type SortValue int

const (
	ASC = "ASC"
	DSC = "DSC"
)

func (s SortValue) String() string {
	return [...]string{"ASC", "DSC"}[s]
}

type SortProps struct {
	Price         string
	Surface       string
	PublishedDate string
}

type PropertiesService struct {
	Log    *slog.Logger
	Client ISSRClient
}

// TODO does the return have to be a pointer?
func NewPropertiesService(log *slog.Logger, client ISSRClient) *PropertiesService {
	return &PropertiesService{
		Log:    log.WithGroup("Properties-Service"),
		Client: client,
	}
}

func (ps *PropertiesService) ListProperties(props SortProps, transactionType string) ([]types.ProcessedListProperty, error) {
	getPropertiesUrl := os.Getenv("GET_PROPERTIES_URL")
	var processedFeatProperties []types.ProcessedListProperty

	constructedUrl := ps.constructGetUrl(transactionType, getPropertiesUrl)

	properties, err := ps.Client.GetPropertiesByTransactionType(constructedUrl)
	if err != nil {
		ps.Log.Error("what's going on here?", "properties", properties)
		return []types.ProcessedListProperty{}, err
	}
	if rootUtils.CheckIfStructIsEmpty(props) {
		if props.Price != "" {
			sortProperties(properties, func(a, b lite.ListPropertiesByTransactionTypeRow) bool {
				if props.Price == ASC {
					return a.Price < b.Price
				} else {
					return a.Price > b.Price
				}
			})
		} else if props.Surface != "" {
			sortProperties(properties, func(a, b lite.ListPropertiesByTransactionTypeRow) bool {
				if props.Price == ASC {
					return a.PropertySurface < b.PropertySurface
				} else {
					return a.PropertySurface > b.PropertySurface
				}
			})
		} else if props.PublishedDate != "" {
			sortProperties(properties, func(a, b lite.ListPropertiesByTransactionTypeRow) bool {
				if props.Price == ASC {
					return a.CreatedAt.Before(b.CreatedAt)
				} else {
					return a.CreatedAt.After(b.CreatedAt)
				}
			})
		}
	}
	for _, featProp := range properties {
		propertyPath := "/" + rootUtils.ReplaceWhiteSpaceWithUnderscore(featProp.Title)
		propertyPath, err = rootUtils.AddParamToUrl(propertyPath, rootUtils.HumanReadableIdQueryKey, featProp.Humanreadableid)
		if err != nil {
			return []types.ProcessedListProperty{}, fmt.Errorf("error trying to create the property path %v", err)
		}

		thumbnailPath := rootUtils.CreateImagePath(featProp.Thumbnail)

		processedFeatProperties = append(processedFeatProperties, types.ProcessedListProperty{
			Title:           featProp.Title,
			TransactionType: featProp.PropertyTransaction,
			Price:           featProp.Price,
			DisplayPrice:    rootUtils.CreateDisplayPrice(featProp.Price),
			PropertyPathUrl: propertyPath,
			CreatedTime:     rootUtils.CreateDisplayCreatedAt(featProp.CreatedAt),
			ThumbnailPath:   thumbnailPath,
		})
	}
	return processedFeatProperties, nil
}

func (ps *PropertiesService) constructGetUrl(transactionType, getUrl string) string {
	var pType = strings.ToLower(transactionType)
	var url = getUrl
	switch pType {
	case Sell:
		url = url + fmt.Sprintf("?transactionType=%s", repositoryTypes.Sell.String())
	case Rent:
		url = url + fmt.Sprintf("?transactionType=%s", repositoryTypes.Rent.String())
	default:
		url = ""
	}

	return url
}

// TODO sorting on the SSR is a bad pattern,
// normally the sorting should be done with database queries
// but it's about trying out how fast golang is
func sortProperties[T any](slice []T, less func(a, b T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
}
