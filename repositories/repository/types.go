package repository

import "errors"

var (
	ErrorNotFound          = errors.New("row doesn't exist")
	ErrorAccessingDatabase = errors.New("error while trying to query the database")
)
