package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func JsonParseAndValidate[T any](ctx *fiber.Ctx, validate *validator.Validate) (*T, JsonError) {
	obj := new(T)
	if err := ctx.BodyParser(obj); err != nil {
		return nil, &HTTPError{Details: err.Error()}
	}
	fmt.Println(obj)
	err := validate.Struct(obj)
	if err != nil {
		verr := err.(validator.ValidationErrors)
		validationError := &ValidationError{
			Fields: nil,
		}
		for _, e := range verr {
			fieldErr := FieldValidationError{
				Name:  e.Field(),
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

func GetToken(ctx *fiber.Ctx) string {
	header := ctx.Get("Authorization")
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
