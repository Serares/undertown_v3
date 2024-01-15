package main

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/login/handlers"
	"github.com/Serares/undertown_v3/services/api/login/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userRepo, err := repository.NewUsersRepository()
	if err != nil {
		log.Error("error on initializing the db", err)
	}
	service := service.NewLoginService(log, userRepo)
	h := handlers.NewLoginHandler(log, &service)
	algnhsa.ListenAndServe(h, nil)
}
