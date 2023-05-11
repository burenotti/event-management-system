package handler

import (
	_ "github.com/burenotti/rtu-it-lab-recruit/docs"
	"github.com/burenotti/rtu-it-lab-recruit/handler/middlewares/auth"
	"github.com/burenotti/rtu-it-lab-recruit/handler/middlewares/logging"
	"github.com/burenotti/rtu-it-lab-recruit/services"
	"github.com/burenotti/rtu-it-lab-recruit/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// HTTPHandler
//
//	@title									API системы управления городскими меропреятиями
//	@version								0.1.0
//	@description							Реализация тестового задания для RTUITLab.
//	@contact.name							Буренин Артём
//	@contact.email							burenotti@gmail.com
//	@license.name							Apache 2.0
//	@license.url							http://www.apache.org/licenses/LICENSE-2.0.html
//	@host									localhost:8000
//	@BasePath								/
//
//	@securitydefinitions.oauth2.password	APIKey
//	@tokenUrl								/auth/sign-in
//	@description							OAuth protects our entity endpoints
type HTTPHandler struct {
	app   *fiber.App
	ucase UseCases
}

type Config struct {
	Name string
}

type UseCases struct {
	services.AuthService
	usecases.EmailSignInUseCase
	usecases.SignUpUseCase
	usecases.OrganizationUseCase
}

func New(logger *logrus.Logger, ucase UseCases, config *Config) *HTTPHandler {
	app := fiber.New(fiber.Config{
		AppName: config.Name,
	})
	app.Use(logging.New(logging.Config{Logger: logger}))
	app.Use(recover.New(recover.Config{}))
	handler := &HTTPHandler{
		ucase: ucase,
		app:   app,
	}
	handler.Mount()
	return handler
}

func (h *HTTPHandler) Mount() {
	authRequired := auth.New(&h.ucase.AuthService)
	h.app.Get("/docs/*", swagger.HandlerDefault)
	auth := h.app.Group("/auth")
	{
		auth.Post("/sign-up", h.SignUp)
		auth.Get("/activate/:token", h.ActivateWithToken)
		auth.Post("/request", h.RequestEmailCode)
		auth.Post("/sign-in", h.SignIn)
	}

	organizations := h.app.Group("/organization", authRequired)
	{
		organizations.Post("/", h.CreateOrganization)
		organizations.Get("/:organization_id", h.GetOrganization)
		organizations.Patch("/:organization_id", h.UpdateOrganization)
		organizations.Delete("/:organization_id", h.DeleteOrganization)
	}
	invites := h.app.Group("/organization/:organization_id/invite", authRequired)
	{
		invites.Post("/", h.InviteToOrganization)
		invites.Get("/", h.ListInvites)
		invites.Post("/:invite_id/accept", h.AcceptInvite)
		invites.Post("/:invite_id/reject", h.RejectInvite)
	}
}

func (h *HTTPHandler) Handler() fasthttp.RequestHandler {
	return h.app.Handler()
}
