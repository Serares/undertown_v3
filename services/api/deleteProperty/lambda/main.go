package main

import (
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/deleteProperty/handlers"
	"github.com/Serares/undertown_v3/services/api/deleteProperty/service"
	"github.com/akrylysov/algnhsa"
	"github.com/joho/godotenv"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Error("error loading the env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3037"
	}

	repo, err := repository.NewPropertiesRepo()
	if err != nil {
		log.Error("error creating the repository")
		return
	}
	service := service.NewDeleteService(log, repo)

	gh := handlers.New(log, service)

	algnhsa.ListenAndServe(gh, nil)
}
