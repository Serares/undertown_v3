package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetSqliteUrl() (string, error) {
	url := os.Getenv("SQLITE_URL")
	port := os.Getenv("SQLITE_PORT")
	db := os.Getenv("SQLITE_DB")
	authToken := os.Getenv("SQLITE_AUTHTOKEN")
	IS_LOCAL := os.Getenv("IS_LOCAL")
	isLocal, err := strconv.ParseBool(IS_LOCAL)
	if err != nil {
		return "", fmt.Errorf("error trying to get the IS_LOCAL env %v", err)
	}
	// if using Turso
	if !isLocal {
		return fmt.Sprintf("libsql://%s.turso.io?authToken=%s", db, authToken), nil
	}

	return fmt.Sprintf("http://%s:%s", url, port), nil
}
