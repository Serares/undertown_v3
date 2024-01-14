package types

import "errors"

var (
	ErrorNotFound          = errors.New("row doesn't exist")
	ErrorAccessingDatabase = errors.New("error while trying to query the database")
)

const (
	PSQL_DRIVER   = "psql"
	SQLITE_DRIVER = "sqlite"
)
