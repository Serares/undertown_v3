package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// TODO this might have to handle the creation of urls for aurora psql also
func CreatePsqlUrl(ctx context.Context, log *slog.Logger) (string, error) {
	dbUser := os.Getenv("PSQL_USER")
	dbPassword := os.Getenv("PSQL_PASSWORD")
	dbName := os.Getenv("PSQL_DB")
	dbHost := os.Getenv("PSQL_HOST")
	dbPort := os.Getenv("PSQL_PORT")
	secretArn := os.Getenv("PSQL_SECRET_ARN")
	if dbUser == "" || dbPassword == "" {
		// log.Info("The context that was passed:", ctx)
		secret, err := getSecret(ctx, secretArn)
		if err != nil {
			log.Error("error trying to access the secret: ", "err", err)
			return "", fmt.Errorf("error trying to access the secret")
		}
		log.Info("Success in getting the secret: ", "secret", secret)
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort), nil
}

// TODO for this to be able to access the secrets manager
// the lambda needs to have IAM authorization
// create the IAM permissions necessary to access the secret resource
func getSecret(ctx context.Context, secretArn string) (interface{}, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	}

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		// handle error
		return nil, err
	}

	var secret interface{}
	if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
		// handle error
		return nil, err
	}

	return secret, err
}
