package handler

import "github.com/gofiber/fiber/v2"

// InviteToOrganization
//
//	@Summary	Invite user to organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Invites
//	@Param		organization_id	path		int					true	"Invite id"
//	@Param		invite			body		model.InviteCreate	true	"Invite"
//	@Success	201				{object}	model.InviteGet
//	@Failure	400				{object}	HTTPError
//	@Failure	500				{object}	HTTPError
//	@Failure	422				{object}	ValidationError
//	@Router		/organization/{organization_id}/invite/ [post]
func (h *HTTPHandler) InviteToOrganization(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// ListInvites
//
//	@Summary	Returns list of invites to organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Invites
//	@Param		organization_id	path		int	true	"Invite id"
//	@Success	201				{object}	[]model.InviteGet
//	@Failure	400				{object}	HTTPError
//	@Failure	500				{object}	HTTPError
//	@Failure	422				{object}	ValidationError
//	@Router		/organization/{organization_id}/invite/ [post]
func (h *HTTPHandler) ListInvites(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// AcceptInvite
//
//	@Summary	Current user accepts invite and joins organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Invites
//	@Param		invite_id		path	int	true	"Invite id"
//	@Param		organization_id	path	int	true	"Organization id"
//	@Success	204
//	@Failure	400	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	ValidationError
//	@Router		/organization/{organization_id}/invite/{invite_id}/accept [post]
func (h *HTTPHandler) AcceptInvite(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// RejectInvite
//
//	@Summary	Current user rejects invite and joins organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Invites
//	@Param		invite_id		path	int	true	"Invite id"
//	@Param		organization_id	path	int	true	"Organization id"
//	@Success	204
//	@Failure	400	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	ValidationError
//	@Router		/organization/{organization_id}/invite/{invite_id}/reject [post]
func (h *HTTPHandler) RejectInvite(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// RemoveMemberFromOrganization
//
//	@Summary	Removes member from organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Members
//	@Param		member_id		path	int	true	"Member id"
//	@Param		organization_id	path	int	true	"Organization id"
//	@Success	204
//	@Failure	400	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	ValidationError
//	@Router		/organization/{organization_id}/member/{member_id}/ [delete]
func (h *HTTPHandler) RemoveMemberFromOrganization(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// UpdatePrivileges
//
//	@Summary	Updates organization member's privileges
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Members
//	@Param		member_id		path	int	true	"Member id"
//	@Param		organization_id	path	int	true	"Organization id"
//	@Success	204
//	@Failure	400	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	ValidationError
//	@Router		/organization/{organization_id}/member/{member_id}/ [put]
func (h *HTTPHandler) UpdatePrivileges(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}

// LeaveOrganization
//
//	@Summary	Current user leaves organization
//	@Security	APIKey
//	@Accept		json
//	@Produce	json
//	@Tags		Members
//	@Param		organization_id	path	int	true	"Invite id"
//	@Success	204
//	@Failure	400	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Failure	422	{object}	ValidationError
//	@Router		/organization/{organization_id}/leave [delete]
func (h *HTTPHandler) LeaveOrganization(ctx *fiber.Ctx) error {
	return (&HTTPError{Details: "NotImplemented"}).AsFiberError(501)
}
