package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/register/types"
	"github.com/joho/godotenv"
)

func setupService(t *testing.T) (*UserService, func()) {
	t.Helper()
	// setup the db
	err := godotenv.Load("../.env.local")
	if err != nil {
		t.Fatalf("Error loading the env file %v", err)
	}

	// setup logger
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// initialize the db
	userRepo, err := repository.NewUsersRepository()
	if err != nil {
		t.Fatal("error connectiong to the db")
	}
	service := NewRegisterService(log, userRepo)
	return &service, func() {
		err := userRepo.CloseDbConnection(context.Background())
		if err != nil {
			t.Errorf("failed to close the db connection for users: %v", err)
		}
	}
}
func TestRegisterService(t *testing.T) {
	cases := []struct {
		title         string
		userData      types.PostUserRequest
		expectedError error
	}{
		{
			title: "Add User",
			userData: types.PostUserRequest{
				Isadmin:  true,
				Issu:     true,
				Email:    "random@email.com",
				Password: "random",
			},
			expectedError: nil,
		},
		{
			title: "Email Exists",
			userData: types.PostUserRequest{
				Isadmin:  true,
				Issu:     true,
				Email:    "random@email.com",
				Password: "random",
			},
			expectedError: types.ErrorEmailAlreadyExists,
		},
	}

	userService, cleanup := setupService(t)
	defer cleanup()

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			err := userService.PersistUser(context.Background(), &tc.userData)
			if err != nil {
				if !errors.Is(err, types.ErrorEmailAlreadyExists) {
					t.Fatalf("expected %v and got %v", types.ErrorEmailAlreadyExists, err)
				}
			}
		})
	}

	// Remove users from db
	for _, tc := range cases {
		ctx := context.Background()
		user, err := userService.UserRepo.GetByEmail(ctx, tc.userData.Email)
		if err != nil {
			if !errors.Is(err, types.ErrorNotFound) {
				t.Fatalf("error searching up mock users %v", err)
			}
			continue
		}
		// remove user
		err = userService.UserRepo.Delete(ctx, user.ID)
		if err != nil {
			t.Fatalf("error cleaning up mock users %v", err)
		}
	}
}
