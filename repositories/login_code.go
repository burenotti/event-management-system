package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/leporo/sqlf"
	"strconv"
	"time"
)

var (
	ErrCodeNotFound = errors.New("code not found")
	ErrCodeInvalid  = errors.New("code is invalid")
)

type LoginCodeRepository struct {
	db      DatabaseWrapper
	CodeTTL time.Duration
}

func (r *LoginCodeRepository) CreateLoginCode(ctx context.Context, userId int64, code string) error {
	if len(code) != 4 {
		return fmt.Errorf("%w: code must be of length 4", ErrCodeInvalid)
	}

	if _, err := strconv.Atoi(code); err != nil {
		return fmt.Errorf("%w: code must be a valid number", err)
	}
	_, err := sqlf.InsertInto("login_code").
		Set("user_id", userId).
		Set("code", code).
		Set("is_used", false).
		Set("expires_at", time.Now().UTC().Add(r.CodeTTL)).
		ExecAndClose(ctx, r.db)

	return err
}

func (r *LoginCodeRepository) MarkCodeUsed(ctx context.Context, userId int64, code string) error {
	res, err := sqlf.Update("login_code").
		Where("user_id = ?", userId).
		Where("code = ?", code).
		Where("is_used = false").
		Where("now() < expires_at").
		Set("is_used", true).
		ExecAndClose(ctx, r.db)

	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrCodeNotFound
	}
	return nil
}
