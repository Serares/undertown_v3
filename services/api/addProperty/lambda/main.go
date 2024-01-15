package main

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/addProperty/handler"
	"github.com/Serares/undertown_v3/services/api/addProperty/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbRepo, err := repository.NewPropertiesRepo()
	ss := service.NewSubmitService(log, dbRepo)
	if err != nil {
		log.Error("error trying to connect to db %v", err)
	}
	addPropertyHandler := handler.New(log, ss)

	algnhsa.ListenAndServe(addPropertyHandler, nil)
}
