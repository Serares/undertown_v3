package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	"github.com/Serares/undertown_v3/services/api/register/handlers"
	"github.com/Serares/undertown_v3/services/api/register/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbUrl, err := utils.CreatePsqlUrl(context.Background(), log)
	if err != nil {
		log.Error("error trying to access the db secret")
	}
	userRepo, err := repository.NewUsersRepository(dbUrl)
	if err != nil {
		log.Error("error on initializing the db")
		os.Exit(1)
	}
	service := service.NewRegisterService(log, userRepo)
	h := handlers.NewRegisterHandler(log, service)
	algnhsa.ListenAndServe(h, nil)
}
