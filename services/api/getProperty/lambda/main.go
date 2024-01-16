package main

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/getProperty/handlers"
	"github.com/Serares/undertown_v3/services/api/getProperty/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbRepo, err := repository.NewPropertiesRepo()
	gps := service.NewGetPropertyService(log, dbRepo)
	if err != nil {
		log.Error("error trying to connect to db %v", err)
	}
	getPropertyHandler := handlers.New(log, gps)

	algnhsa.ListenAndServe(getPropertyHandler, nil)
}
