package handler

import (
	_ "github.com/burenotti/rtu-it-lab-recruit/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/valyala/fasthttp"
)

// HTTPHandler
//
//	@title			API системы управления городскими меропреятиями
//	@version		0.1.0
//	@description	Реализация тестового задания для RTUITLab.
//	@contact.name	Буренин Артём
//	@contact.email	burenotti@gmail.com
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8000
//	@BasePath		/
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
	h.app.Get("/docs/*", swagger.HandlerDefault)
	auth := h.app.Group("/auth")
	{
		auth.Post("/sign-up", h.SignUp)
		auth.Get("/activate/:token", h.ActivateWithToken)
		auth.Post("/request", h.RequestEmailCode)
		auth.Post("/sign-in", h.SignIn)
	}
}
