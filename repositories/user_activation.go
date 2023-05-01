package repositories

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

const (
	DefaultIssuer = "RTUITLab"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

type UserActivationRepository struct {
	PrivateKey *rsa.PrivateKey
	TokenTTL   time.Duration
}

func NewUserActivationRepositoryFromPem(privateKeyPem []byte) (*UserActivationRepository, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPem)
	if err != nil {
		return nil, err
	}
	return &UserActivationRepository{
		PrivateKey: privateKey,
		TokenTTL:   0,
	}, nil
}

func NewUserActivationRepositoryFromFile(privateKeyPath string) (*UserActivationRepository, error) {
	data, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return NewUserActivationRepositoryFromPem(data)
}

func (r *UserActivationRepository) CreateActivationToken(_ context.Context, userId int64) (*model.ActivationToken, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(r.TokenTTL)

	claims := model.ActivationTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    DefaultIssuer,
			Subject:   fmt.Sprintf("%d", userId),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(r.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &model.ActivationToken{
		Token:     tokenString,
		UserId:    userId,
		ExpiresAt: expiresAt,
	}, nil
}

func (r *UserActivationRepository) ValidateActivationToken(_ context.Context, token string) (*model.ActivationToken, error) {
	jwtToken, err := jwt.Parse(token, r.selectKey)
	if errors.Is(err, ErrInvalidToken) {
		return nil, err
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, fmt.Errorf("%w: provieded token is not a jwt token", ErrInvalidToken)
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		exp, _ := jwtToken.Claims.GetExpirationTime()
		return nil, fmt.Errorf("%w: token exired at %s", ErrTokenExpired, exp.String())
	} else if err != nil {
		return nil, err
	}

	exp, _ := jwtToken.Claims.GetExpirationTime()
	sub, _ := jwtToken.Claims.GetSubject()
	var userId int64
	_, err = fmt.Sscanf(sub, "%d", &userId)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid token subject: %s", ErrInvalidToken, sub)
	}
	return &model.ActivationToken{
		Token:     token,
		UserId:    userId,
		ExpiresAt: exp.Time.UTC(),
	}, nil
}

func (r *UserActivationRepository) selectKey(token *jwt.Token) (interface{}, error) {
	if token.Method != jwt.SigningMethodRS256 {
		return nil, fmt.Errorf("%w: %s method is not supported", ErrInvalidToken, token.Method.Alg())
	}

	return &r.PrivateKey.PublicKey, nil
}
