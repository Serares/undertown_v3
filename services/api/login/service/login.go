package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	repositoryTypes "github.com/Serares/undertown_v3/repositories/repository/types"
	"github.com/Serares/undertown_v3/services/api/login/types"
	"github.com/Serares/undertown_v3/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	Log      *slog.Logger
	UserRepo repository.IUsersRepository
}

func NewLoginService(log *slog.Logger, userrepo repository.IUsersRepository) LoginService {
	return LoginService{
		Log:      log,
		UserRepo: userrepo,
	}
}

// returns the base64 jwt if success
func (ls *LoginService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	user, err := ls.searchUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("error searching for user email %w", err)
	}
	// compare if passwords match
	// base64 decode the hashedpassword also
	hashedPassword, err := base64.RawStdEncoding.Strict().DecodeString(user.Passwordhash)
	if err != nil {
		return "", fmt.Errorf("error decoding hashed password %w", err)
	}
	err = ls.comparePasswords(password, hashedPassword)
	if err != nil {
		if errors.Is(types.ErrorWrongPassword, err) {
			// TODO try to unify this message with the login handler also
			return "", fmt.Errorf("email or password is wrong")
		}
		ls.Log.Error("error while trying to decrypt the passwords", "error", err)
		return "", fmt.Errorf("error while trying to login %w", err)
	}

	// generate the JWT
	token, err := ls.signAndGenerateToken(user.Email, user.ID, user.Isadmin)
	if err != nil {
		return "", fmt.Errorf("error while generating token %w", err)
	}
	// TODO convert token to base64
	// tokenInBase64 := base64.RawStdEncoding.EncodeToString([]byte(token))
	return token, nil
}

func (ls *LoginService) searchUserByEmail(ctx context.Context, email string) (*lite.User, error) {
	user, err := ls.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositoryTypes.ErrorNotFound) {
			// have to return error in case the user is not existent
			return nil, fmt.Errorf("user email is not registered")
		}
		return nil, err
	}
	return user, nil
}

func (ls *LoginService) comparePasswords(plainPassword string, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	if err != nil {
		return types.ErrorWrongPassword
	}
	return nil
}
func (ls *LoginService) signAndGenerateToken(email string, userId string, isadmin bool) (string, error) {
	claims := utils.JWTClaims{
		Email:   email,
		UserId:  userId,
		Isadmin: isadmin,
		IsSsr:   true,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
