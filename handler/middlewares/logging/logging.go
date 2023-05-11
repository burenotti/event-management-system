package logging

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Logger *logrus.Logger
}

type Middleware struct {
	Logger *logrus.Logger
}

func (m *Middleware) Call(ctx *fiber.Ctx) error {
	log := m.Logger.
		WithField("time", time.Now().UTC()).
		WithField("method", ctx.Method()).
		WithField("path", ctx.Path())

	err := ctx.Next()
	if err != nil {
		log = m.Logger.WithError(err)
	}
	log.WithField("status", ctx.Response().StatusCode()).
		WithField("execution_time", time.Now().UTC()).
		Info()

	return err
}

func New(cfg Config) fiber.Handler {
	m := Middleware{
		Logger: cfg.Logger,
	}
	return m.Call
}
