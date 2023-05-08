package handler

import (
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

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

	u, jerr := JsonParseAndValidate[model.UserCreate](ctx, validate)
	if jerr != nil {
		return jerr.AsFiberError(fiber.StatusUnprocessableEntity)
	}
	user, err := h.ucase.SignUp(ctx.Context(), u)
	if err != nil {
		return WrapError(err)
	}

	return ReturnJson(ctx, user)
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
	token := ctx.Params("token")
	if token == "" {
		return NewHTTPError("token is required path parameter").AsFiberError(fiber.StatusBadRequest)
	}
	if err := h.ucase.ActivateWithToken(ctx.Context(), token); err != nil {
		return WrapError(err)
	}
	ctx.Status(fiber.StatusNoContent)
	return nil
}

// RequestEmailCode
//
//	@Tags		Auth
//	@Summary	Requests sending one time password to users email
//	@Accept		x-www-form-urlencoded
//	@Produce	json
//
//	@Param		email	formData	string	true	"Request data"
//
//	@Success	204
//	@Failure	500	{object}	HTTPError
//	@Failure	400	{object}	HTTPError
//	@Router		/auth/request [post]
func (h *HTTPHandler) RequestEmailCode(ctx *fiber.Ctx) error {
	req, jerr := JsonParseAndValidate[model.CodeRequest](ctx, validate)
	if jerr != nil {
		return jerr.AsFiberError(422)
	}

	if err := h.ucase.RequestCode(ctx.Context(), req.Email); err != nil {
		return WrapError(err)
	}
	return nil
}

// SignIn
//
//	@Tags		Auth
//	@Summary	Signs user in using sent in email one time password
//	@Accept		x-www-form-urlencoded
//	@Produce	json
//
//	@Param		username	formData	string	true	"Email"
//	@Param		password	formData	string	true	"One time code"
//
//	@Success	200			{object}	model.Token
//	@Failure	500			{object}	HTTPError
//	@Failure	422			{object}	HTTPError
//	@Failure	401			{object}	HTTPError
//	@Router		/auth/sign-in [post]
func (h *HTTPHandler) SignIn(ctx *fiber.Ctx) error {
	req, jerr := JsonParseAndValidate[model.AuthCredentials](ctx, validate)
	//fmt.Println(req.Email)
	if jerr != nil {
		return jerr.AsFiberError(422)
	}
	token, err := h.ucase.SignIn(ctx.Context(), req)

	if err != nil {
		return WrapError(err)
	}

	return ReturnJson(ctx, model.NewAccessToken(token))
}
