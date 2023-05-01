package model

type OrganizationCreate struct {
	OrganizationID int64   `json:"organization_id"`
	Name           string  `json:"name"`
	Address        *string `json:"address,omitempty"`
	ContactEmail   *string `json:"contact_email,omitempty"`
	ContactPhone   *string `json:"contact_phone,omitempty"`
}

type Organization struct {
	OrganizationID int64   `json:"organization_id"`
	Name           string  `json:"name"`
	Address        *string `json:"address,omitempty"`
	ContactEmail   *string `json:"contact_email,omitempty"`
	ContactPhone   *string `json:"contact_phone,omitempty"`
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
