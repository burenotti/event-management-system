package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type ActivationToken struct {
	Token     string
	UserId    int64
	ExpiresAt time.Time
}

type ActivationTokenClaims struct {
	jwt.RegisteredClaims
}
