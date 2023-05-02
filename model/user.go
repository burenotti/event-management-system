package model

type UserCreate struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
}

type UserGet struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	IsActive   bool   `json:"is_active"`
}

type User struct {
	UserID     int64  `json:"user_id" faker:"-"`
	FirstName  string `json:"first_name" faker:"first_name"`
	LastName   string `json:"last_name" faker:"last_name"`
	MiddleName string `json:"middle_name" faker:"-"`
	Email      string `json:"email" faker:"email"`
	IsActive   bool   `json:"is_active" faker:"-"`
}
