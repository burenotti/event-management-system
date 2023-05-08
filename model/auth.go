package model

import "github.com/golang-jwt/jwt/v5"

type AuthCredentials struct {
	Email string `json:"username" form:"username" validate:"required,email,gte=1" example:"johndoe@example.com"`
	Code  string `json:"password" form:"password" validate:"required,numeric,len=4" example:"6666"`
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

type CodeRequest struct {
	Email string `json:"email" validate:"required,email" example:"johndoe@example.com"`
}

type Token struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

func NewAccessToken(token string) *Token {
	return &Token{
		Type:        "Bearer",
		AccessToken: token,
	}
}
