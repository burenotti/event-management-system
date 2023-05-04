package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func JsonParseAndValidate[T any](validate *validator.Validate, body []byte) (*T, JsonError) {
	obj := new(T)
	if err := json.Unmarshal(body, obj); err != nil {
		return nil, &HTTPError{Details: "Cannot parse json object"}
	}
	err := validate.Struct(obj)
	if err != nil {
		verr := err.(*validator.ValidationErrors)
		validationError := &ValidationError{
			Fields: nil,
		}
		for _, e := range *verr {
			fieldErr := FieldValidationError{
				Name:  e.StructField(),
				Error: e.Error(),
			}
			validationError.Fields = append(validationError.Fields, fieldErr)
		}
		return nil, validationError
	}
	return obj, nil
}

func ReturnJson(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.JSON(obj); err != nil {
		return NewHTTPError(err.Error()).AsFiberError(fiber.StatusInternalServerError)
	}
	return nil
}
