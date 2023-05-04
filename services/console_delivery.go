package services

import (
	"context"
	"github.com/sirupsen/logrus"
)

type ConsoleDelivery struct {
	Logger *logrus.Logger
}

func (c *ConsoleDelivery) SendCode(_ context.Context, email string, code string) error {
	c.Logger.
		WithField("email", email).
		WithField("code", code).
		Infof("ConsoleDelivery new login code")

	return nil
}

func (c *ConsoleDelivery) SendToken(_ context.Context, email string, token string) error {
	c.Logger.
		WithField("email", email).
		WithField("token", token).
		Infof("ConsoleDelivery new activation token")

	return nil
}
