package auth

import (
	"context"
	"errors"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/burenotti/rtu-it-lab-recruit/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

var (
	ErrTokenIsNotProvided = errors.New("bearer token is not provided")
)

const UserCtxKey = "user_payload"

type Middleware struct {
	authService *services.AuthService
}

func New(auth *services.AuthService) fiber.Handler {
	m := Middleware{authService: auth}
	return m.Call
}

func (m *Middleware) Call(ctx *fiber.Ctx) error {
	scheme, token, err := GetToken(ctx)
	if err != nil || strings.ToLower(scheme) != "bearer" {
		ctx.Set("WWW-Authenticate", "Bearer")
		return fiber.NewError(fiber.StatusUnauthorized, "Not authenticated")
	}
	payload, err := m.authService.ValidateToken(ctx.Context(), token)
	if err != nil {
		ctx.Set("WWW-Authenticate", "Bearer")
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	ctx.SetUserContext(context.WithValue(ctx.Context(), UserCtxKey, payload))
	return ctx.Next()
}

func GetToken(ctx *fiber.Ctx) (scheme string, token string, err error) {
	header := ctx.Get("Authorization")

	if header == "" {
		return "", "", ErrTokenIsNotProvided
	}

	parts := strings.Split(header, " ")

	if len(parts) != 2 {
		return "", "", ErrTokenIsNotProvided
	}

	return parts[0], parts[1], err
}

func GetAuth(ctx *fiber.Ctx) (*model.AuthPayload, bool) {
	auth, ok := ctx.Context().Value(UserCtxKey).(*model.AuthPayload)
	return auth, ok
}
