package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/psql"
	"github.com/Serares/undertown_v3/services/api/register/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type UserService struct {
	Log      *slog.Logger
	UserRepo repository.IUsersRepository
}

func NewRegisterService(log *slog.Logger, userRepo repository.IUsersRepository) UserService {
	return UserService{
		Log:      log,
		UserRepo: userRepo,
	}
}

func generateSalt(length int32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (us *UserService) checkUserExists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := us.UserRepo.Get(ctx, id)
	// only return false if the SqlNoRows error is returned by the repo
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return false, nil
		}
		us.Log.Error("%s -- %v", types.ErrorCheckingIfUserExists, err)
		return false, err
	}
	return true, nil
}

func (us *UserService) checkIfEmailAlreadyExists(ctx context.Context, email string) (bool, error) {
	_, err := us.UserRepo.GetByEmail(ctx, email)
	// only return false if the SqlNoRows error is returned by the repo
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return false, nil
		}
		us.Log.Error("%s -- %v", types.ErrorCheckingIfEmailExists, err)
		return false, err
	}
	return true, nil
}

// The hash is stored as base64
func (us *UserService) createPasswordHash(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 1, 32)
}

func (us *UserService) PersistUser(ctx context.Context, user *types.PostUserRequest) error {
	// create the user params to be persisted
	userParams := psql.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Isadmin:   user.Isadmin,
		Issu:      user.Issu,
		Email:     user.Email,
	}
	isEmailExists, err := us.checkIfEmailAlreadyExists(ctx, userParams.Email)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if isEmailExists {
		return types.ErrorEmailAlreadyExists
	}
	// This might be redundant
	// but it's checking for id clashes
	isExist, err := us.checkUserExists(ctx, userParams.ID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	// Create the password hash after checking if the user exist to not waste resources
	if !isExist {
		salt, err := generateSalt(8)
		fmt.Printf("the salt in string is: %s\n", string(salt))
		if err != nil {
			// TODO try to retry salt creation a couple of times
			return fmt.Errorf("error creating salt for password")
		}
		hash := us.createPasswordHash(user.Password, salt)
		if err != nil {
			return fmt.Errorf("%s -- %v", types.ErrorHashingPassword, err)
		}
		fmt.Printf("The generated hash as byte: %s \n", hash)
		userParams.Passwordhash = base64.RawStdEncoding.EncodeToString(hash)
		// password hash is stored as base64 also
		userParams.Passwordsalt = base64.RawStdEncoding.EncodeToString(salt)

		fmt.Printf("Decoding the salt: %s", string(salt))
		decodedSalt, err := base64.RawStdEncoding.Strict().DecodeString(userParams.Passwordsalt)
		if err != nil {
			fmt.Printf("error decoding salt: %v", err)
		}
		fmt.Printf("the decoded salt: %v\n", decodedSalt)

		fmt.Printf("Password hash: %s \n", userParams.Passwordhash)
		decoded, err := base64.StdEncoding.Strict().DecodeString(userParams.Passwordhash)
		if err != nil {
			fmt.Printf("error decoding the hash: %v\n", err)
		}
		fmt.Printf("password hash decoded: %s \n", string(decoded))
	}
	// persist the user
	return us.UserRepo.Add(ctx, userParams)
}
