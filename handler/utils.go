package handler

import (
	"errors"
	"github.com/burenotti/rtu-it-lab-recruit/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

func JsonParseAndValidate[T any](ctx *fiber.Ctx, validate *validator.Validate) (*T, JsonError) {
	obj := new(T)
	if err := ctx.BodyParser(obj); err != nil {
		return nil, &HTTPError{Details: err.Error()}
	}
	err := validate.Struct(obj)
	if err != nil {
		verr := err.(validator.ValidationErrors)
		validationError := &ValidationError{
			Fields: nil,
		}
		for _, e := range verr {
			field, _ := reflect.TypeOf(*obj).FieldByName(e.StructField())
			fieldName := strings.Split(field.Tag.Get("json"), ",")[0]
			fieldErr := FieldValidationError{
				Name:  fieldName,
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

func UnwrapAtomicError(err error) error {
	if err == nil {
		return nil
	}
	if atErr, ok := err.(*repositories.AtomicError); ok {
		if atErr.TransactionError != nil && !errors.Is(atErr.InnerError, atErr.TransactionError) {
			return atErr.TransactionError
		} else if atErr.InnerError != nil {
			return atErr.InnerError
		} else {
			return nil
		}
	} else {
		return err
	}
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	err = UnwrapAtomicError(err)

	httpError := NewHTTPError(err.Error())
	if errors.Is(err, repositories.ErrUserNotFound) {
		return httpError.AsFiberError(fiber.StatusBadRequest)
	} else if errors.Is(err, repositories.ErrUserExists) {
		return httpError.AsFiberError(fiber.StatusBadRequest)
	} else if errors.Is(err, repositories.ErrCodeInvalid) {
		return httpError.AsFiberError(fiber.StatusForbidden)
	} else if errors.Is(err, repositories.ErrMemberNotFound) {
		return httpError.AsFiberError(fiber.StatusBadRequest)
	} else if errors.Is(err, repositories.ErrLogicError) {
		return httpError.AsFiberError(fiber.StatusBadRequest)
	} else if errors.Is(err, repositories.ErrInvalidToken) {
		return httpError.AsFiberError(fiber.StatusBadRequest)
	} else if errors.Is(err, repositories.ErrCodeNotFound) {
		return httpError.AsFiberError(fiber.StatusForbidden)
	} else if errors.Is(err, repositories.ErrCodeInvalid) {
		return httpError.AsFiberError(fiber.StatusBadRequest)
	} else {
		return httpError.AsFiberError(fiber.StatusInternalServerError)
	}
}
