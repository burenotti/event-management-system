package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/burenotti/rtu-it-lab-recruit/repositories"
	"github.com/sirupsen/logrus"
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
	SendActivationToken(ctx context.Context, user *model.User, token string) error
}

type SignUpUseCase struct {
	UserRepo       UserStorage
	ActivationRepo UserActivator
	Delivery       TokenDelivery
	Logger         *logrus.Logger
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

	u.Logger.
		WithField("user_id", user.UserID).
		Infof("Sending activation email to %d", user.UserID)
	if err = u.Delivery.SendActivationToken(ctx, user, token.Token); err != nil {
		u.Logger.
			WithField("user_id", user.UserID).
			WithError(err).
			Info("Sending email activation failed: %v", err)
	}

	return user, nil
}

func (u *SignUpUseCase) ActivateWithToken(ctx context.Context, token string) error {
	t, err := u.ActivationRepo.ValidateActivationToken(ctx, token)
	if err != nil {
		return err
	}
	u.Logger.WithField("user_id", t.UserId).Info("Activating user %d", t.UserId)

	_, err = u.UserRepo.Update(ctx, t.UserId, repositories.UpdatesMap{
		"is_active": true,
	})
	return err
}
