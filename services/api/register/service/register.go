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
	"golang.org/x/crypto/bcrypt"
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
func (us *UserService) createPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
		hash, err := us.createPasswordHash(user.Password)
		if err != nil {
			return fmt.Errorf("%s -- %v", types.ErrorHashingPassword, err)
		}
		userParams.Passwordhash = base64.RawStdEncoding.EncodeToString(hash)
	}
	// persist the user
	return us.UserRepo.Add(ctx, userParams)
}
