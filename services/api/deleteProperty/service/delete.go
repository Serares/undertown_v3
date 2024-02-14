package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Serares/undertown_v3/repositories/repository"
)

type DeleteService struct {
	Log            *slog.Logger
	PropertiesRepo *repository.Properties
}

func NewDeleteService(log *slog.Logger, propertiesRepo *repository.Properties) DeleteService {
	return DeleteService{
		Log:            log,
		PropertiesRepo: propertiesRepo,
	}
}

func (s DeleteService) DeleteProperty(ctx context.Context, id, humanReadableId string) error {

	prop, err := s.PropertiesRepo.GetById(ctx, "", humanReadableId)
	if err != nil {
		s.Log.Error(
			"error trying to get the property data before deletion",
			"error", err,
		)
	}

	// Dispatch the delete images SQS event and then delete the property

	err = s.PropertiesRepo.DeleteByHumanReadableId(ctx, humanReadableId)
	if err != nil {
		s.Log.Error("error when trying to delete the property", "error", err)
		return fmt.Errorf("error trying to delete the property with human readable id %w", err)
	}

	// ❗TODO dispatch the SQS request to delete processed images

	return nil
}
