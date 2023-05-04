package services

import (
	"context"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/sirupsen/logrus"
)

type ConsoleDelivery struct {
	Logger *logrus.Logger
}

func (c *ConsoleDelivery) SendCode(_ context.Context, user *model.User, code string) error {
	c.Logger.
		WithField("email", user.Email).
		WithField("code", code).
		Infof("ConsoleDelivery new login code")

	return nil
}

func (c *ConsoleDelivery) SendActivationToken(_ context.Context, user *model.User, token string) error {
	c.Logger.
		WithField("email", user.Email).
		WithField("token", token).
		Infof("ConsoleDelivery new activation token")

	return nil
}
