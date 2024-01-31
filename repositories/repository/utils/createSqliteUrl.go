package utils

import (
	"fmt"
	"os"
)

func CreateSqliteUrl() (string, error) {
	host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	protocol := os.Getenv("DB_PROTOCOL")
	dbLocal := os.Getenv("DB_LOCAL")
	dbName := os.Getenv("DB_NAME")
	authToken := os.Getenv("TURSO_DB_TOKEN")
	if dbLocal == "true" {
		return host, nil
	}

	// return a turso url
	return fmt.Sprintf("%s://%s.%s?authToken=%s", protocol, dbName, host, authToken), nil
}
