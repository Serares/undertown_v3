package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type GetPropertiesService struct {
	Log            *slog.Logger
	PropertiesRepo *repository.Properties
}

func NewPropertiesService(log *slog.Logger, propertiesRepo *repository.Properties) GetPropertiesService {
	return GetPropertiesService{
		Log:            log,
		PropertiesRepo: propertiesRepo,
	}
}

func (s GetPropertiesService) ListFeaturedProperties(ctx context.Context) (*[]lite.ListFeaturedPropertiesRow, error) {
	properties, err := s.PropertiesRepo.ListFeatured(ctx)
	if err != nil {
		return nil, fmt.Errorf("error trying to get featured properties %w", err)
	}

	return &properties, nil
}

func (s GetPropertiesService) ListProperties(ctx context.Context) (*[]lite.Property, error) {
	properties, err := s.PropertiesRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error trying to get properties %w", err)
	}

	return &properties, nil
}
