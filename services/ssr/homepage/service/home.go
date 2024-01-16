package service

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type IHomeClient interface {
	ListFeaturedProperties(string) ([]lite.ListFeaturedPropertiesRow, error)
}

type HomeService struct {
	Log    *slog.Logger
	client IHomeClient
}

func NewHomeService(log *slog.Logger, client IHomeClient) *HomeService {
	return &HomeService{
		Log:    log,
		client: client,
	}
}

func (hs *HomeService) ListProperties() ([]lite.ListFeaturedPropertiesRow, error) {
	getPropertiesUrl := os.Getenv("GET_PROPERTIES_URL")

	return hs.client.ListFeaturedProperties(strings.Join([]string{getPropertiesUrl, "isFeatured=true"}, "?"))
}
