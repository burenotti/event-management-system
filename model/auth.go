package model

type AuthCredentials struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}
