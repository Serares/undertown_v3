package service

import "github.com/Serares/undertown_v3/repositories/repository/lite"

type ISSRClient interface {
	ListFeaturedProperties(string) ([]lite.ListFeaturedPropertiesRow, error)
	GetProperty(url string) (lite.Property, error)
	GetProperties(url string) ([]lite.Property, error)
	GetPropertiesByTransactionType(url string) ([]lite.ListPropertiesByTransactionTypeRow, error)
}
