package model

type UserCreate struct {
	FirstName  string `json:"first_name" example:"John"`
	LastName   string `json:"last_name" example:"Doe"`
	MiddleName string `json:"middle_name" example:"Jr."`
	Email      string `json:"email" example:"johndoe@example.com"`
}

type UserGet struct {
	UserID     int64  `json:"user_id" example:"1"`
	FirstName  string `json:"first_name" example:"John"`
	LastName   string `json:"last_name" example:"Doe"`
	MiddleName string `json:"middle_name" example:"Jr."`
	Email      string `json:"email" example:"johndoe@example.com"`
	IsActive   bool   `json:"is_active" example:"true"`
}

type User struct {
	UserID     int64  `json:"user_id" faker:"-"`
	FirstName  string `json:"first_name" faker:"first_name"`
	LastName   string `json:"last_name" faker:"last_name"`
	MiddleName string `json:"middle_name" faker:"-"`
	Email      string `json:"email" faker:"email"`
	IsActive   bool   `json:"is_active" faker:"-"`
}
