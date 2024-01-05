package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Serares/api/register/types"
	"github.com/Serares/repository"
	"github.com/Serares/repository/psql"
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

func (us *UserService) createPasswordHash(password string, salt []byte) (string, error) {
	hash := argon2.IDKey([]byte(password), salt, 2, 64*1024, 1, 32)
	return base64.RawStdEncoding.EncodeToString(hash), nil
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
		return fmt.Errorf("%s", types.ErrorEmailAlreadyExists)
	}
	isExist, err := us.checkUserExists(ctx, userParams.ID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	// Create the password hash after checking if the user exist to not waste resources
	if !isExist {
		salt, err := generateSalt(8)
		if err != nil {
			// TODO try to retry salt creation a couple of times
			return fmt.Errorf("error creating salt for password")
		}
		hash, err := us.createPasswordHash(user.Password, salt)
		if err != nil {
			return fmt.Errorf("%s -- %v", types.ErrorHashingPassword, err)
		}
		userParams.Passwordhash = hash
		userParams.Passwordsalt = base64.RawStdEncoding.EncodeToString(salt)
	}
	// persist the user
	return us.UserRepo.Add(ctx, userParams)
}
