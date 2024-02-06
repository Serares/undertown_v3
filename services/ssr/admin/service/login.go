package service

import (
	"log/slog"
	"os"
)

type LoginService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewLoginService(log *slog.Logger, client ISSRAdminClient) *LoginService {
	return &LoginService{
		Log:    log.WithGroup("Login Service"),
		Client: client,
	}
}

func (s *LoginService) Login(email, password string) (string, error) {
	// TODO
	// run some validations here if needed
	// get the token from cookie

	loginUrl := os.Getenv("LOGIN_URL")
	return s.Client.Login(email, password, loginUrl)
}
