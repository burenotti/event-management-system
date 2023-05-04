package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/burenotti/rtu-it-lab-recruit/repositories"
)

var (
	ErrTokenSendFailed = errors.New("activation token send failed")
)

type UserStorage interface {
	Create(ctx context.Context, u *model.UserCreate) (*model.User, error)
	GetById(ctx context.Context, userId int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, userId int64, update map[string]interface{}) (*model.User, error)
	Delete(ctx context.Context, userId int64) error
}

type UserActivator interface {
	CreateActivationToken(_ context.Context, userId int64) (*model.ActivationToken, error)
	ValidateActivationToken(_ context.Context, token string) (*model.ActivationToken, error)
}

type TokenDelivery interface {
	SendToken(ctx context.Context, email string, token string) error
}

type SignUpUseCase struct {
	UserRepo       UserStorage
	ActivationRepo UserActivator
	Delivery       TokenDelivery
}

func (u *SignUpUseCase) SignUp(ctx context.Context, create *model.UserCreate) (*model.User, error) {
	user, err := u.UserRepo.Create(ctx, create)
	if err != nil {
		return nil, err
	}

	token, err := u.ActivationRepo.CreateActivationToken(ctx, user.UserID)
	if err != nil {
		return user, fmt.Errorf("%w: can't create activation token", ErrTokenSendFailed)
	}

	if err = u.Delivery.SendToken(ctx, user.Email, token.Token); err != nil {
		return user, fmt.Errorf("%w: token Delivery failed", ErrTokenSendFailed)
	}
	return user, nil
}

func (u *SignUpUseCase) ActivateWithToken(ctx context.Context, token string) error {
	t, err := u.ActivationRepo.ValidateActivationToken(ctx, token)
	if err != nil {
		return err
	}

	_, err = u.UserRepo.Update(ctx, t.UserId, repositories.UpdatesMap{
		"is_active": true,
	})
	return err
}
