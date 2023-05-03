package handler

import "github.com/gofiber/fiber/v2"

// SignUp
//
//	@Tags		Auth
//	@Summary	Creates new user that should be activated with email
//	@Accept		json
//	@Produce	json
//	@Param		user	body		model.UserCreate	true	"User info"
//	@Success	201		{object}	model.UserGet
//	@Failure	422		{object}	HTTPError
//	@Failure	500		{object}	HTTPError
//	@Failure	400		{object}	HTTPError
//	@Router		/auth/sign-up [post]
func (h *HTTPHandler) SignUp(ctx *fiber.Ctx) error {
	panic("Unimplemented")

}

// ActivateWithToken
//
//	@Tags		Auth
//	@Summary	Activates user with token sent in email
//	@Accept		json
//	@Produce	json
//
//	@Param		token	path	string	true	"Activation token"
//
//	@Success	204
//	@Failure	500	{object}	HTTPError
//	@Failure	400	{object}	HTTPError
//	@Router		/auth/activate/{token} [get]
func (h *HTTPHandler) ActivateWithToken(ctx *fiber.Ctx) error {
	panic("Unimplemented")
}

// RequestEmailCode
//
//	@Tags			Auth
//	@Summary		Requests sending one time password to users email
//	@Description	In real API activate method should be POST, because it changes system state.
//	@Description	But in this case it was made GET to make possible activation by clicking on link in email message without any additional code.
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body	model.CodeRequest	true	"Request data"
//
//	@Success		204
//	@Failure		500	{object}	HTTPError
//	@Failure		400	{object}	HTTPError
//	@Router			/auth/request [post]
func (h *HTTPHandler) RequestEmailCode(ctx *fiber.Ctx) error {
	panic("Unimplemented")
}

// SignIn
//
//	@Tags		Auth
//	@Summary	Signs user in using sent in email one time password
//	@Accept		json
//	@Produce	json
//
//	@Param		request	body	model.AuthCredentials	true	"Credentials"
//
//	@Success	200
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	HTTPError
//	@Failure	401	{object}	HTTPError
//	@Router		/auth/sign-in [post]
func (h *HTTPHandler) SignIn(ctx *fiber.Ctx) error {
	panic("Unimplemented")
}
