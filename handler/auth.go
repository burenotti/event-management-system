package handler

import (
	"fmt"
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
		return NewHTTPError(err.Error()).AsFiberError(fiber.StatusBadRequest)
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
		return NewHTTPError("token is required path parameter").AsFiberError(400)
	}
	if err := h.ucase.ActivateWithToken(ctx.Context(), token); err != nil {
		return NewHTTPError(err.Error()).AsFiberError(400)
	}
	ctx.Status(204)
	return nil
}

// RequestEmailCode
//
//	@Tags		Auth
//	@Summary	Requests sending one time password to users email
//	@Accept		json
//	@Produce	json
//
//	@Param		request	body	model.CodeRequest	true	"Request data"
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
		return NewHTTPError(err.Error()).AsFiberError(400)
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
//	@Param		username	query	string	true	"Email"
//	@Param		password	query	string	true	"One time code"
//
//	@Success	200
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	HTTPError
//	@Failure	401	{object}	HTTPError
//	@Router		/auth/sign-in [post]
func (h *HTTPHandler) SignIn(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Get("Content-Type"))
	req, jerr := JsonParseAndValidate[model.AuthCredentials](ctx, validate)
	//fmt.Println(req.Email)
	if jerr != nil {
		return jerr.AsFiberError(422)
	}
	token, err := h.ucase.SignIn(ctx.Context(), req)

	if err != nil {
		return NewHTTPError(err.Error()).AsFiberError(400)
	}

	return ReturnJson(ctx, model.NewAccessToken(token))
}
