package model

import "github.com/golang-jwt/jwt/v5"

type AuthCredentials struct {
	UserId int64  `json:"user_id"`
	Code   string `json:"code"`
}

type AuthTokenClaims struct {
	jwt.RegisteredClaims
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthPayload struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
