package service

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/utils"
)

type DeleteService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewDeleteService(log *slog.Logger, client ISSRAdminClient) *DeleteService {
	return &DeleteService{
		Log:    log.WithGroup("Delete Service"),
		Client: client,
	}
}

func (ds *DeleteService) Delete(humanReadableId, id, authToken string) error {
	deleteUrl := os.Getenv("DELETE_PROPERTY_URL")
	var err error
	var url string
	if humanReadableId != "" {
		url, err = utils.AddParamToUrl(deleteUrl, utils.HumanReadableIdQueryKey, humanReadableId)
		if err != nil {
			return err
		}
	}

	if id != "" {
		url, err = utils.AddParamToUrl(deleteUrl, utils.HumanReadableIdQueryKey, id)
		if err != nil {
			return err
		}
	}

	return ds.Client.DeleteProperty(url, authToken)
}
