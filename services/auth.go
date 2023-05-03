package services

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)

type AuthService struct {
	TokenTTL   time.Duration
	privateKey *rsa.PrivateKey
}

func (s *AuthService) CreateToken(_ context.Context, user *model.User) (string, error) {
	now := jwt.NewNumericDate(time.Now().UTC())
	expiresAt := jwt.NewNumericDate(now.Add(s.TokenTTL))
	claims := model.AuthTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   fmt.Sprintf("%d", user.UserID),
			ExpiresAt: expiresAt,
			NotBefore: now,
			IssuedAt:  now,
		},
		UserID:    user.UserID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(s.privateKey)
	return token, err
}

func (s *AuthService) ValidateToken(_ context.Context, tokenString string) (*model.AuthPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return &s.privateKey.PublicKey, nil
	})
	payload := &model.AuthPayload{}
	claims, _ := token.Claims.(jwt.MapClaims)
	data, _ := json.Marshal(claims)
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
