package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	repositoryTypes "github.com/Serares/undertown_v3/repositories/repository/types"
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

func (us *UserService) checkUserExists(ctx context.Context, id string) (bool, error) {
	// derefference a struct
	_, err := us.UserRepo.Get(ctx, id)
	// only return false if the SqlNoRows error is returned by the repo
	if err != nil {
		if errors.Is(err, repositoryTypes.ErrorNotFound) {
			return false, nil
		}
		us.Log.Error("error checking if user exists", "type", types.ErrorCheckingIfUserExists, "error", err)
		return false, err
	}
	return true, nil
}

func (us *UserService) checkIfEmailAlreadyExists(ctx context.Context, email string) (bool, error) {
	_, err := us.UserRepo.GetByEmail(ctx, email)
	// only return false if the SqlNoRows error is returned by the repo
	if err != nil {
		us.Log.Error("error checking for email", err)
		if errors.Is(err, repositoryTypes.ErrorNotFound) {
			return false, nil
		}
		us.Log.Error("error checking if email exists", "type", types.ErrorCheckingIfEmailExists, "error", err)
		return false, err
	}
	return true, nil
}

// The hash is stored as base64
func (us *UserService) createPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (us *UserService) PersistUser(ctx context.Context, user *types.PostUserRequest) error {
	us.Log.Info("persist user method", "user", user)
	// create the user params to be persisted
	userParams := lite.CreateUserParams{
		ID:        uuid.New().String(),
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
		us.Log.Info("Creating password hash")
		hash, err := us.createPasswordHash(user.Password)
		if err != nil {
			return fmt.Errorf("%s -- %v", types.ErrorHashingPassword, err)
		}
		userParams.Passwordhash = base64.RawStdEncoding.EncodeToString(hash)
	}
	// persist the user
	return us.UserRepo.Add(ctx, userParams)
}
