package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

// /⚠️ TODO do some performance testing here
// check how the go routines are spawned
// view the colStats project

func TestMain(t *testing.M) {
	os.Exit(t.Run())
}

// todo use the local admin submit
// to create a movk sqs message body and use it for testing
func TestPersist(t *testing.T) {
	err := godotenv.Load("../.env.local")
	if err != nil {
		t.Errorf("error loading the .env file")
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo, err := repository.NewPropertiesRepo()
	if err != nil {
		t.Errorf("error trying to initialize the repo %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Errorf(
			"error trying to load the lambda context %v",
			err,
		)
	}
	sqsClient := sqs.NewFromConfig(cfg)

	service := NewPUService(
		log,
		repo,
		sqsClient,
	)
	t.Run("Processing the sqs body and trying to persist the property", func(t *testing.T) {
		mockedSqsBody, err := os.ReadFile("../testdata/mockSqsError.json")
		if err != nil {
			t.Errorf("error trying to read the json mock for sqs body %v", err)
		}

		err = service.Persist(context.TODO(), string(mockedSqsBody), "whatever")
		if err != nil {
			t.Errorf("error trying to process the json mock for sqs body %v", err)
		}
	})
}
