package service

import (
	"io"
	"log/slog"
	"os"
)

type SubmitService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewSubmitService(log *slog.Logger, client ISSRAdminClient) SubmitService {
	return SubmitService{
		Log:    log,
		Client: client,
	}
}

func (s SubmitService) Submit(body io.Reader, authToken string) error {
	// TODO
	// run some validations here if needed
	// get the token from cookie

	submitUrl := os.Getenv("SUBMIT_PROPERTY_URL")
	return s.Client.AddProperty(body, submitUrl, authToken)
}
