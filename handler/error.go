package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type JsonError interface {
	Json() string
	AsFiberError(status int) error
}

type HTTPError struct {
	Details string `json:"details"`
}

func NewHTTPError(details string) *HTTPError {
	return &HTTPError{Details: details}
}

func (e *HTTPError) AsFiberError(status int) error {
	return fiber.NewError(status, e.Json())
}

func (e *HTTPError) Json() string {
	data, _ := json.Marshal(e)
	return string(data)
}

type FieldValidationError struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

type ValidationError struct {
	Fields []FieldValidationError `json:"fields"`
}

func (e *ValidationError) AsFiberError(status int) error {
	return fiber.NewError(status, e.Json())
}

func (e *ValidationError) Json() string {
	data, _ := json.Marshal(e)
	return string(data)
}
