package main

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/getProperties/handlers"
	"github.com/Serares/undertown_v3/services/api/getProperties/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbRepo, err := repository.NewPropertiesRepo()
	ss := service.NewPropertiesService(log, dbRepo)
	if err != nil {
		log.Error("error trying to connect to db %v", err)
	}
	addPropertyHandler := handlers.New(log, ss)

	algnhsa.ListenAndServe(addPropertyHandler, nil)
}
