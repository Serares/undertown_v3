package service

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type EditService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewEditService(log *slog.Logger, client ISSRAdminClient) *EditService {
	return &EditService{
		Log:    log,
		Client: client,
	}
}

// TODO right now the Submit() method from the SubmitService is used to edit the property
// func (es *EditService) Post(body, humanReadableId, authToken string) error {
// }

func (es *EditService) Get(humanReadableId, authToken string) (lite.Property, error) {
	getPropertyUrl := os.Getenv("GET_PROPERTY_URL")
	return es.Client.GetProperty(getPropertyUrl, humanReadableId, authToken)
}
