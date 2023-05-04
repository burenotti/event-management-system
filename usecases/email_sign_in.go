package usecases

import (
	"context"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/burenotti/rtu-it-lab-recruit/repositories"
	"math/rand"
)

type StorageTransactioner interface {
	Atomic(ctx context.Context, f repositories.AtomicFunc) error
}

type LoginCodeStorage interface {
	CreateLoginCode(ctx context.Context, userId int64, code string) error
	MarkCodeUsed(ctx context.Context, userId int64, code string) error
}

type CodeDelivery interface {
	SendCode(ctx context.Context, email string, code string) error
}

type AuthService interface {
	CreateToken(ctx context.Context, user *model.User) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (*model.AuthPayload, error)
}

type EmailSignInUseCase struct {
	UserStore      UserStorage
	LoginCodeStore LoginCodeStorage
	Delivery       CodeDelivery
	Transactioner  StorageTransactioner
	Auth           AuthService
}

func (s *EmailSignInUseCase) RequestCode(ctx context.Context, userEmail string) error {
	return s.Transactioner.Atomic(ctx, func(ctx context.Context) error {
		user, err := s.UserStore.GetByEmail(ctx, userEmail)
		if err != nil {
			return err
		}

		code := GenerateActivationCode()
		if err = s.LoginCodeStore.CreateLoginCode(ctx, user.UserID, code); err != nil {
			return err
		}

		text := fmt.Sprintf("Your activation code is: %s", code)
		err = s.Delivery.SendCode(ctx, user.Email, text)
		return err
	})
}

func (s *EmailSignInUseCase) SignIn(ctx context.Context, creds *model.AuthCredentials) (string, error) {
	var token string
	err := s.Transactioner.Atomic(ctx, func(ctx context.Context) error {
		user, err := s.UserStore.GetByEmail(ctx, creds.Email)
		if err != nil {
			return err
		}

		if err := s.LoginCodeStore.MarkCodeUsed(ctx, user.UserID, creds.Code); err != nil {
			return err
		}

		token, err = s.Auth.CreateToken(ctx, user)
		return err
	})
	return token, err
}

func GenerateActivationCode() string {
	intCode := rand.Int() % 10000
	return fmt.Sprintf("%04d", intCode)
}
