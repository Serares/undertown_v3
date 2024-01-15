package main

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/register/handlers"
	"github.com/Serares/undertown_v3/services/api/register/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	userRepo, err := repository.NewUsersRepository()
	if err != nil {
		log.Error("error on initializing the db", "error:", err)
		os.Exit(1)
	}
	service := service.NewRegisterService(log, userRepo)
	h := handlers.NewRegisterHandler(log, service)
	algnhsa.ListenAndServe(h, nil)
}
