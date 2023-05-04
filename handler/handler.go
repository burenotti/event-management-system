package handler

import (
	_ "github.com/burenotti/rtu-it-lab-recruit/docs"
	"github.com/burenotti/rtu-it-lab-recruit/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	app   *fiber.App
	ucase UseCases
}

type Config struct {
	Name string
}

type UseCases struct {
	usecases.EmailSignInUseCase
	usecases.SignUpUseCase
}

func New(ucase UseCases, config *Config) *HTTPHandler {
	app := fiber.New(fiber.Config{
		AppName: config.Name,
	})
	app.Use(recover.New())
	handler := &HTTPHandler{
		ucase: ucase,
		app:   app,
	}
	handler.Mount()
	return handler
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

func (h *HTTPHandler) Handler() fasthttp.RequestHandler {
	return h.app.Handler()
}
