package handler

import "github.com/gofiber/fiber/v2"

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
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
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
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
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
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}
