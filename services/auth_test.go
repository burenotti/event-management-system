package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mrand "math/rand"
	"testing"
	"time"
)

func TestAuthService_CreateToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	s := AuthService{
		TokenTTL:   24 * time.Hour,
		PrivateKey: key,
	}
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		u := model.User{}
		u.UserID = int64(mrand.Int() % 1000)
		err = faker.FakeData(&u)
		t.Run(fmt.Sprintf("%s %s", u.FirstName, u.LastName), func(t *testing.T) {
			require.NoError(t, err)
			ss, err := s.CreateToken(ctx, &u)
			assert.NoError(t, err, "should create token without errors")
			token, err := jwt.Parse(ss, func(token *jwt.Token) (interface{}, error) {
				return key, nil
			})
			claims, ok := token.Claims.(jwt.MapClaims)
			assert.True(t, ok, "should cast claims to MapClaims")
			assert.EqualValues(t, int(u.UserID), int(claims["user_id"].(float64)))
			assert.EqualValues(t, u.FirstName, claims["first_name"])
			assert.EqualValues(t, u.LastName, claims["last_name"])
			assert.EqualValues(t, u.Email, claims["email"])
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	s := AuthService{
		TokenTTL:   24 * time.Hour,
		PrivateKey: key,
	}
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		u := model.User{}
		u.UserID = int64(mrand.Int() % 1000)
		err = faker.FakeData(&u)
		t.Run(fmt.Sprintf("%s %s", u.FirstName, u.LastName), func(t *testing.T) {
			token, err := s.CreateToken(ctx, &u)
			require.NoError(t, err)
			p, err := s.ValidateToken(ctx, token)
			assert.NoError(t, err)
			assert.Equal(t, u.UserID, p.UserID)
			assert.Equal(t, u.Email, p.Email)
			assert.Equal(t, u.FirstName, p.FirstName)
			assert.Equal(t, u.LastName, p.LastName)
		})
	}
}
