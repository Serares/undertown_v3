package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/deleteProperty/handlers"
	"github.com/Serares/undertown_v3/services/api/deleteProperty/service"
	"github.com/akrylysov/algnhsa"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Error(
			"error trying to load the lambda context",
			"error", err,
		)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	repo, err := repository.NewPropertiesRepo()
	if err != nil {
		log.Error("error creating the repository")
		return
	}
	service := service.NewDeleteService(log, repo, sqsClient)

	gh := handlers.New(log, service)

	algnhsa.ListenAndServe(gh, nil)
}
