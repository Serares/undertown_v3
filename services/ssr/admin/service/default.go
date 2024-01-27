package service

import "github.com/Serares/undertown_v3/repositories/repository/lite"

type ISSRClient interface {
	AddProperty(params lite.AddPropertyParams) ([]lite.ListFeaturedPropertiesRow, error)
	Login(email, password string) (string, error)
	UpdateProperty() error // todo have to implement the update lambda
	GetProperty(id string) (lite.Property, error)
}
