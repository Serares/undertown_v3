package types

import "errors"

var (
	ErrorCheckingIfUserExists  = errors.New("error trying to search for user in db")
	ErrorPersistingUser        = errors.New("error trying to persist user in db")
	ErrorHashingPassword       = errors.New("error hashing the password")
	ErrorCheckingIfEmailExists = errors.New("error trying to search for user email")
	ErrorEmailAlreadyExists    = errors.New("email already exists")
	ErrorNotFound              = errors.New("not found")
)

type PostUserRequest struct {
	Isadmin  bool   `json:"isadmin"`
	Issu     bool   `json:"issu"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
