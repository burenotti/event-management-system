package model

type Invite struct {
	InviteID       int64  `json:"invite_id" example:"1"`
	UserId         int64  `json:"user_id" example:"2"`
	UserEmail      string `json:"user_email" example:"johndoe@example.com"`
	ToOrganization int64  `json:"to_organization" example:"1"`
	FromUser       int64  `json:"from_user" example:"1"`
	Status         string `json:"status" enums:"sent,accepted,rejected"`
}

type InviteGet struct {
	InviteID       int64  `json:"invite_id" example:"1"`
	UserId         int64  `json:"user_id" example:"2"`
	UserEmail      string `json:"user_email" example:"johndoe@example.com"`
	ToOrganization int64  `json:"to_organization" example:"1"`
	FromUser       int64  `json:"from_user" example:"1"`
	Status         string `json:"status" enums:"sent,accepted,rejected"`
}

type InviteCreate struct {
	UserEmail      string `json:"user_email" example:"johndoe@example.com"`
	ToOrganization int64  `json:"to_organization" example:"1"`
	FromUser       int64  `json:"from_user" example:"1"`
}
