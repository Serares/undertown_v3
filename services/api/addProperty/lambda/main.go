package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/addProperty/handler"
	"github.com/Serares/undertown_v3/services/api/addProperty/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbUrl := createPsqlUrl()
	dbRepo, err := repository.NewPropertiesRepo(dbUrl)
	ss := service.NewSubmitService(log, dbRepo)
	if err != nil {
		log.Error("error trying to connect to db %v", err)
	}
	addPropertyHandler := handler.New(log, ss)

	algnhsa.ListenAndServe(addPropertyHandler, nil)
}

func createPsqlUrl() string {
	dbUser := os.Getenv("PSQL_USER")
	dbPassword := os.Getenv("PSQL_PASSWORD")
	dbName := os.Getenv("PSQL_NAME")
	dbHost := os.Getenv("PSQL_HOST")
	dbPort := os.Getenv("PSQL_PORT")

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
}
