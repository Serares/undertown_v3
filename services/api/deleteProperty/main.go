package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/deleteProperty/handlers"
	"github.com/Serares/undertown_v3/services/api/deleteProperty/service"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      gh,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
