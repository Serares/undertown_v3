package service

import (
	"context"
	"log/slog"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type GetPropertyService struct {
	Log            *slog.Logger
	PropertiesRepo *repository.Properties
}

func NewGetPropertyService(log *slog.Logger, repo *repository.Properties) GetPropertyService {
	return GetPropertyService{
		Log:            log,
		PropertiesRepo: repo,
	}
}

// ❗
// Todo the property has to be stripped of sensitive columns, like user_id and the id
func (gp GetPropertyService) GetPropertyByHumanReadableId(ctx context.Context, humanReadableid string) (lite.Property, error) {
	// validate the human readable id
	return gp.PropertiesRepo.GetById(ctx, "", humanReadableid)
}

func (gp GetPropertyService) GetPropertyById(ctx context.Context, id string) (lite.Property, error) {
	// validate the id
	return gp.PropertiesRepo.GetById(ctx, id, "")
}
