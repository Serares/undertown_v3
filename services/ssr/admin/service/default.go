package service

import (
	"io"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type ISSRAdminClient interface {
	AddProperty(body io.Reader, url, authToken, contentType string) error
	Login(email, password, url string) (string, error)
	List(url, authToken string) ([]lite.Property, error)
	GetProperty(url, humanReadableId, authToken string) (lite.Property, error)
	EditProperty(body io.Reader, url, authToken, contentType string) error
}
