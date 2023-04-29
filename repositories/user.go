package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/leporo/sqlf"
	"github.com/pkg/errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrLogicError   = errors.New("logic error")
)

type UserRepository struct {
	db DatabaseWrapper
}

func NewUserRepository(db DatabaseWrapper) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *model.UserCreate) (*model.User, error) {
	var userId int64 = 0
	err := sqlf.InsertInto("users").
		Set("first_name", u.FirstName).
		Set("last_name", u.LastName).
		Set("middle_name", u.MiddleName).
		Set("email", u.Email).
		Set("is_active", false).
		Returning("user_id").To(&userId).
		QueryRow(ctx, r.db)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		UserID:     userId,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		MiddleName: u.MiddleName,
		Email:      u.Email,
		IsActive:   false,
	}
	return user, nil
}

func (r *UserRepository) GetById(ctx context.Context, userId int64) (*model.User, error) {
	u := &model.User{}
	err := sqlf.From("users").
		Select("user_id").To(&u.UserID).
		Select("first_name").To(&u.FirstName).
		Select("last_name").To(&u.LastName).
		Select("middle_name").To(&u.MiddleName).
		Select("is_active").To(&u.IsActive).
		Select("email").To(&u.Email).
		Where("user_id = ?", userId).
		QueryRow(ctx, r.db)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%v: user with provided id does not exist", ErrUserNotFound)
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, userId int64, update map[string]interface{}) (*model.User, error) {

	// Set of columns, which could be updated
	fields := map[string]struct{}{
		"first_name": {}, "last_name": {}, "middle_name": {},
		"email": {}, "is_active": {},
	}

	u := &model.User{}

	query := sqlf.Update("users").
		Where("user_id = ? ", userId).
		Returning("user_id").To(&u.UserID).
		Returning("first_name").To(&u.FirstName).
		Returning("last_name").To(&u.LastName).
		Returning("middle_name").To(&u.MiddleName).
		Returning("email").To(&u.Email).
		Returning("is_active").To(&u.IsActive)

	for key, upd := range update {
		if _, ok := fields[key]; !ok {
			return nil, fmt.Errorf("%v: could not update field '%s'", ErrLogicError, upd)
		}
		query = query.Set(key, upd)
	}

	err := query.QueryRow(ctx, r.db)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%v: user with provided id does not exist", ErrUserNotFound)
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Delete(ctx context.Context, userId int64) error {
	res, err := sqlf.DeleteFrom("users").
		Where("user_id = ?", userId).
		Exec(ctx, r.db)

	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("%v: user with provided id does not exist", ErrUserNotFound)
	}
	return nil
}
