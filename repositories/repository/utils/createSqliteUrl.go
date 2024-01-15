package utils

import (
	"fmt"
	"os"
)

func CreateSqliteUrl() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	protocol := os.Getenv("DB_PROTOCOL")
	isLocal := os.Getenv("IS_LOCAL")
	dbName := os.Getenv("DB_NAME")
	authToken := os.Getenv("TURSO_DB_TOKEN")
	if isLocal == "true" {
		return fmt.Sprintf("%s://%s:%s", protocol, host, port), nil
	}

	// return a turso url
	return fmt.Sprintf("%s://%s.%s?authToken=%s", protocol, dbName, host, authToken), nil
}
