package handler

import (
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/handler/middlewares/auth"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/gofiber/fiber/v2"
)

var ErrTokenIsNotProvided = (&HTTPError{Details: "token is not provided"}).AsFiberError(401)

// CreateOrganization
//
//	@Summary	Creates a new organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Organizations
//	@Param		request	body		model.OrganizationCreate	true	"Organization Info"
//	@Success	201		{object}	model.OrganizationGet
//	@Failure	400		{object}	HTTPError
//	@Failure	500		{object}	HTTPError
//	@Failure	422		{object}	ValidationError
//	@Router		/organization/ [post]
func (h *HTTPHandler) CreateOrganization(ctx *fiber.Ctx) error {

	org, jerr := JsonParseAndValidate[model.OrganizationCreate](ctx, validate)
	if jerr != nil {
		return jerr.AsFiberError(422)
	}

	user, _ := auth.GetAuth(ctx)

	createdOrg, err := h.ucase.CreateOrganization(ctx.Context(), user.UserID, org)
	if err != nil {
		return NewHTTPError(err.Error()).AsFiberError(400)
	}
	ctx.Status(201)
	return ReturnJson(ctx, createdOrg)
}

// GetOrganization
//
//	@Summary	Returns an information about organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Organizations
//	@Param		organization_id	path		int	true	"Organization id"
//	@Success	200				{object}	model.OrganizationGet
//	@Failure	400				{object}	HTTPError
//	@Failure	500				{object}	HTTPError
//	@Router		/organization/{organization_id} [get]
func (h *HTTPHandler) GetOrganization(ctx *fiber.Ctx) error {
	orgId, err := getOrganizationId(ctx)
	if err != nil {
		return err
	}

	org, err := h.ucase.OrganizationUseCase.GetOrganization(ctx.Context(), orgId)
	if err != nil {
		return WrapError(err)
	}
	ctx.Status(fiber.StatusOK)
	return ReturnJson(ctx, org)
}

// UpdateOrganization
//
//	@Summary	Updates organization information
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Organizations
//	@Param		organization_id	path	int							true	"Organization id"
//	@Param		updates			body	model.OrganizationUpdate	true	"Fields that will be updated"
//	@Success	204
//	@Failure	400	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Failure	500	{object}	ValidationError
//	@Router		/organization/{organization_id} [patch]
func (h *HTTPHandler) UpdateOrganization(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// DeleteOrganization
//
//	@Summary	Deletes organization by id
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Organizations
//	@Param		organization_id	path		int	true	"Organization id"
//	@Failure	400				{object}	HTTPError
//	@Failure	500				{object}	HTTPError
//	@Router		/organization/{organization_id} [delete]
func (h *HTTPHandler) DeleteOrganization(ctx *fiber.Ctx) error {

	user, _ := auth.GetAuth(ctx)

	orgId, err := getOrganizationId(ctx)
	if err != nil {
		return err
	}

	err = h.ucase.OrganizationUseCase.DeleteOrganization(ctx.Context(), user, orgId)
	ctx.Status(fiber.StatusOK)
	return WrapError(err)
}

func getOrganizationId(ctx *fiber.Ctx) (int64, error) {
	orgIdRaw := ctx.Params("organization_id")
	if orgIdRaw == "" {
		return 0, NewHTTPError("organization_id is required path parameter").
			AsFiberError(422)
	}
	var orgId int64
	if _, err := fmt.Sscanf(orgIdRaw, "%d", &orgId); err != nil {
		return 0, NewHTTPError("organization_id must be a number").AsFiberError(422)
	}
	return orgId, nil
}
