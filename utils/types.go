package utils

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Email   string `json:"email"`
	UserId  string `json:"userId"`
	Isadmin bool   `json:"isadmin"`
	IsSsr   bool   `json:"isssr"`
	jwt.RegisteredClaims
}
