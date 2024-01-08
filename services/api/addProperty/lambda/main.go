package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	"github.com/Serares/undertown_v3/services/api/addProperty/handler"
	"github.com/Serares/undertown_v3/services/api/addProperty/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbUrl, err := utils.CreatePsqlUrl(context.Background(), log)
	if err != nil {
		log.Error("error on creating the connection string")
	}
	dbRepo, err := repository.NewPropertiesRepo(dbUrl)
	ss := service.NewSubmitService(log, dbRepo)
	if err != nil {
		log.Error("error trying to connect to db %v", err)
	}
	addPropertyHandler := handler.New(log, ss)

	algnhsa.ListenAndServe(addPropertyHandler, nil)
}
