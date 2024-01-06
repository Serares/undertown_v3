package types

import "errors"

var (
	// don't show this message to the users
	ErrorWrongPassword = errors.New("wrong password")
)
