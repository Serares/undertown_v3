package service

import (
	"io"
)

type ISSRAdminClient interface {
	AddProperty(body io.Reader, url, authToken string) error
	Login(email, password, url string) (string, error)
}
