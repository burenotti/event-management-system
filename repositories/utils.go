package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
	"testing"
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

func CreateRandomUser(ctx context.Context, db DatabaseWrapper, t *testing.T) *model.User {
	query := `
		INSERT INTO 
    		users (first_name, last_name, middle_name, email, is_active)
		VALUES ($1, $2, $3, $4, $5) RETURNING user_id
	`
	u := model.User{}
	err := faker.FakeData(&u)
	require.NoError(t, err, "faker generate error")
	row := db.QueryRowContext(ctx, query, u.FirstName, u.LastName, u.MiddleName, u.Email, true)
	err = row.Scan(&u.UserID)
	require.NoError(t, err, "should create user without errors")
	return &u
}
