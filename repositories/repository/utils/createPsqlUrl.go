package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// This is exactly what's stored in the secrets manager variable as json
// the port should be the default pg port
type DBSecretValue struct {
	Host                string `json:"host"`
	Password            string `json:"password"`
	Dbname              string `json:"dbname"`
	Username            string `json:"username"`
	DbClusterIdentified string `json:"dbClusterIdentifier"`
}

// TODO this might have to handle the creation of urls for aurora psql also
func CreatePsqlUrl(ctx context.Context) (string, error) {
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
			return "", fmt.Errorf("error trying to access the secret")
		}
		dbHost = secret.Host
		dbPassword = secret.Password
		dbName = secret.Dbname
		dbUser = secret.Username
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort), nil
}

// TODO for this to be able to access the secrets manager
// the lambda needs to have IAM authorization
// create the IAM permissions necessary to access the secret resource
func getSecret(ctx context.Context, secretArn string) (*DBSecretValue, error) {
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

	var secret DBSecretValue
	if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
		// handle error
		return nil, err
	}

	return &secret, err
}
