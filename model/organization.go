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
