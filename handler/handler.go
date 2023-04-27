package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// HTTPHandler
// @title API системы управления городскими меропреятиями
// @version 0.1.0
// @description Реализация тестового задания для RTUITLab.
// @contact.name Буренин Артём
// @contact.email burenotti@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host app:80
// @BasePath /
type HTTPHandler struct {
	app *fiber.App
}

type Config struct {
	Name string
}

func New(config *Config) fasthttp.RequestHandler {
	handler := &HTTPHandler{}
	handler.app = fiber.New(fiber.Config{
		AppName: config.Name,
	})
	handler.Mount()
	return handler.app.Handler()
}

func (h *HTTPHandler) Mount() {
	h.app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"docs": "http://docs/",
		})
	})
}
