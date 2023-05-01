package repositories

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
)

var (
	ErrUpdatesValidationError = errors.New("updates validation error")
)

func getViolatedConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.ConstraintName
	} else {
		return ""
	}
}

type UpdatesMap map[string]interface{}

type UpdatesValidator struct {
	allowedFields map[string]struct{}
}

func NewUpdatesValidator(allowedFields []string) *UpdatesValidator {
	v := UpdatesValidator{
		allowedFields: make(map[string]struct{}, len(allowedFields)),
	}
	for _, field := range allowedFields {
		v.allowedFields[field] = struct{}{}
	}
	return &v
}

func (v *UpdatesValidator) Validate(updates UpdatesMap) error {
	for field := range updates {
		if _, ok := v.allowedFields[field]; !ok {
			return fmt.Errorf("%w: field '%s' is not allowed", ErrUpdatesValidationError, field)
		}
	}
	return nil
}
