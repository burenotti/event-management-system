package repositories

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func TestNewUserActivationRepositoryFromPem(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
	require.NoError(t, err, "should generate key without errors")
	bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	})

	repo, err := NewUserActivationRepositoryFromPem(pemEncoded)
	assert.NoError(t, err, "should correctly parse private key")
	assert.True(t, privateKey.Equal(repo.PrivateKey))
}

func TestNewUserActivationRepositoryFromFile(t *testing.T) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
	require.NoError(t, err, "should generate key without errors")
	bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	})
	filepath := path.Join(t.TempDir(), "test-private.pem")
	err = ioutils.AtomicWriteFile(filepath, pemEncoded, os.FileMode(777))
	require.NoError(t, err, "should correctly write private key into file")

	repo, err := NewUserActivationRepositoryFromFile(filepath)
	assert.NoError(t, err, "should correctly parse private key")
	assert.True(t, privateKey.Equal(repo.PrivateKey))
}

func TestUserActivationRepository_CreateActivationToken(t *testing.T) {
	ctx := context.Background()
	privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
	require.NoError(t, err, "should generate key without errors")

	repo := UserActivationRepository{
		PrivateKey: privateKey,
		TokenTTL:   24 * time.Hour,
	}
	var userId int64 = 1
	token, err := repo.CreateActivationToken(ctx, userId)
	assert.NoError(t, err, "should correctly create activation token")
	expectedExpiration := time.Now().Add(repo.TokenTTL).UTC().Unix()
	expiration := token.ExpiresAt.UTC().Unix()
	// Since CreateActivationToken uses time.Now() inside,
	// expected expiration time could differ a bit from actual
	assert.InEpsilon(t, expectedExpiration, expiration, 60, "expiration should be calculated correctly")
	jwtToken, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	assert.NoError(t, err, "token should be parsed correctly")
	sub, _ := jwtToken.Claims.GetSubject()
	assert.Equal(t, "1", sub, "user id should be marshaled to jwt as subject")
}

func TestUserActivationRepository_ValidateActivationToken(t *testing.T) {
	ctx := context.Background()
	privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
	require.NoError(t, err, "should generate key without errors")

	repo := UserActivationRepository{
		PrivateKey: privateKey,
		TokenTTL:   24 * time.Hour,
	}
	now := time.Now().UTC()
	expires := now.Add(repo.TokenTTL)
	claims := model.ActivationTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    DefaultIssuer,
			Subject:   "1",
			ExpiresAt: jwt.NewNumericDate(expires),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	require.NoError(t, err, "token should be created without errors")
	token, err := repo.ValidateActivationToken(ctx, tokenString)
	assert.NoError(t, err, "validation of correct token should not lead to error")
	assert.Equal(t, int64(1), token.UserId)
	assert.Equal(t, tokenString, token.Token)
	assert.Equal(t, expires.Truncate(time.Second), token.ExpiresAt)

	// Test malformed token
	_, err = repo.ValidateActivationToken(ctx, "not a token")
	assert.ErrorIs(t, err, ErrInvalidToken)

	// Test wrong signing method
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(privateKey)
	require.NoError(t, err, "token should be created without errors")
	_, err = repo.ValidateActivationToken(ctx, tokenString)
	assert.ErrorIs(t, err, ErrInvalidToken)

	// Test expired token
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC())
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	require.NoError(t, err, "token should be created without errors")
	_, err = repo.ValidateActivationToken(ctx, tokenString)
	assert.ErrorIs(t, err, ErrTokenExpired)

}
