package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInjectContext(t *testing.T) {
	ctx := context.Background()
	tx := &sqlx.Tx{}
	txCtx := injectTx(ctx, tx)
	assert.Equal(t, tx, txCtx.Value(txKey{}), "should be equal")
}

func TestExtractContext(t *testing.T) {
	ctx := context.Background()
	tx := &sqlx.Tx{}
	txCtx := injectTx(ctx, tx)
	actualTx := extractTx(txCtx)
	assert.Equal(t, tx, actualTx, "should be equal")
}

func TestExtractContext_NoTransaction(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, extractTx(ctx), "should return nil, if no tx in context")
}

func TestNewDatabase(t *testing.T) {
	db := &sqlx.DB{}
	database := NewDatabase(db)
	assert.Equal(t, db, database.db)
}

func TestDatabase_WithContext_NoTx(t *testing.T) {
	ctx := context.Background()
	db := &sqlx.DB{}
	database := NewDatabase(db)
	actual, ok := database.WithContext(ctx).(*sqlx.DB)
	assert.True(t, ok)
	assert.Equal(t, db, actual)
}

func TestDatabase_WithContext(t *testing.T) {
	ctx := context.Background()
	tx := &sqlx.Tx{}
	ctx = injectTx(ctx, tx)
	db := &sqlx.DB{}
	database := NewDatabase(db)
	actual, ok := database.WithContext(ctx).(*sqlx.Tx)
	assert.True(t, ok)
	assert.Equal(t, tx, actual)
}

func TestDatabase_Atomic_Ok(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "postgres")
	require.NoError(t, err, "could not create sql mock")

	database := NewDatabase(dbx)
	mock.ExpectBegin()
	mock.ExpectCommit()
	err = database.Atomic(ctx, func(ctx context.Context) error {
		assert.NotNil(t, extractTx(ctx), "context should contain transaction")
		return nil
	})
	assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should met")
	assert.Nil(t, err, "closure doesn't return any errors, so atomic shouldn't")
}

func TestDatabase_Atomic_ClosureReturnsErr(t *testing.T) {
	expectedErr := errors.New("test Error")
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "postgres")
	require.NoError(t, err, "could not create sql mock")

	database := NewDatabase(dbx)
	mock.ExpectBegin()
	mock.ExpectRollback()
	err = database.Atomic(ctx, func(ctx context.Context) error {
		assert.NotNil(t, extractTx(ctx), "context should contain transaction")
		return expectedErr
	})
	assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should met")
	assert.ErrorAsf(t, err, &expectedErr, "should return wrapped test error")
}

func TestDatabase_Atomic_ClosurePanics(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "postgres")
	require.NoError(t, err, "could not create sql mock")

	database := NewDatabase(dbx)
	mock.ExpectBegin()
	mock.ExpectRollback()
	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "test", r, "should ret")
			}
		}()
		err = database.Atomic(ctx, func(ctx context.Context) error {
			panic("test")
		})
		assert.Fail(t, "this line should not be reached because of panic")
	}()

	assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should met")
}

func TestDatabase_Atomic_CantBegin(t *testing.T) {
	expectedErr := errors.New("test")
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "postgres")
	require.NoError(t, err, "could not create sql mock")

	mock.ExpectBegin().WillReturnError(expectedErr)
	database := NewDatabase(dbx)
	err = database.Atomic(ctx, func(ctx context.Context) error {
		assert.Fail(t, "this line should not be reached because of panic")
		return nil
	})
	assert.ErrorAs(t, err, &expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should met")
}

func TestDatabase_Atomic_CantCommit(t *testing.T) {
	expectedErr := errors.New("test")
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "postgres")
	require.NoError(t, err, "could not create sql mock")

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(expectedErr)
	database := NewDatabase(dbx)
	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.ErrorContains(t, r.(error), "test", "should ret")
			}
		}()
		err = database.Atomic(ctx, func(ctx context.Context) error {
			panic("test")
		})
		assert.Fail(t, "this line should not be reached because of panic")
	}()
	assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should met")
}
