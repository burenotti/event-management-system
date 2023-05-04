package model

type OrganizationCreate struct {
	Name         string  `json:"name" example:"Российский технологический университет МИРЭА"`
	Address      *string `json:"address,omitempty" example:"Г. Москва, Пр-т. Вернадского 78"`
	ContactEmail *string `json:"contact_email,omitempty" example:"contact@mirea.ru"`
	ContactPhone *string `json:"contact_phone,omitempty" example:"74992156565"`
}

type Organization struct {
	OrganizationID int64   `json:"organization_id"`
	Name           string  `json:"name"`
	Address        *string `json:"address,omitempty"`
	ContactEmail   *string `json:"contact_email,omitempty"`
	ContactPhone   *string `json:"contact_phone,omitempty"`
}

type OrganizationGet struct {
	OrganizationID int64   `json:"organization_id" example:"1"`
	Name           string  `json:"name" example:"Российский технологический университет МИРЭА"`
	Address        *string `json:"address,omitempty" example:"Г. Москва, Пр-т. Вернадского 78"`
	ContactEmail   *string `json:"contact_email,omitempty" example:"contact@mirea.ru"`
	ContactPhone   *string `json:"contact_phone,omitempty" example:"74992156565"`
}

type OrganizationUpdate struct {
	Name         *string `json:"name" example:"Российский технологический университет МИРЭА"`
	Address      *string `json:"address,omitempty" example:"Г. Москва, Пр-т. Вернадского 78"`
	ContactEmail *string `json:"contact_email,omitempty" example:"contact@mirea.ru"`
	ContactPhone *string `json:"contact_phone,omitempty" example:"74992156565"`
}

type OrganizationMemberCreate struct {
	UserID int64            `json:"user_id"`
	Can    MemberPrivileges `json:"privileges"`
}

type OrganizationMember struct {
	UserID int64            `json:"user_id"`
	Can    MemberPrivileges `json:"privileges"`
}

type MemberPrivileges struct {
	ViewEvents    bool `json:"view_events"`
	EditEvents    bool `json:"edit_events"`
	ManageMembers bool `json:"manage_members"`
}
