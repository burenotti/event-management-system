package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/leporo/sqlf"
)

func init() {
	sqlf.SetDialect(sqlf.PostgreSQL)
}

type AtomicFunc func(ctx context.Context) error

type DbContext interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
	sqlx.PreparerContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DatabaseWrapper interface {
	DbContext
	Atomic(ctx context.Context, f AtomicFunc) error
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	} else {
		return nil
	}
}

type Database struct {
	db *sqlx.DB
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.withContext(ctx).QueryContext(ctx, query, args...)
}

func (d *Database) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.withContext(ctx).QueryxContext(ctx, query, args...)
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.withContext(ctx).QueryRowContext(ctx, query, args...)
}

func (d *Database) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return d.withContext(ctx).QueryRowxContext(ctx, query, args...)
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.withContext(ctx).ExecContext(ctx, query, args...)
}

func (d *Database) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.withContext(ctx).PrepareContext(ctx, query)
}

func (d *Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.withContext(ctx).SelectContext(ctx, dest, query, args...)
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) withContext(ctx context.Context) DbContext {
	if tx := extractTx(ctx); tx != nil {
		return tx
	} else {
		return d.db
	}
}

func (d *Database) Atomic(ctx context.Context, f AtomicFunc) error {

	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})

	defer func() {
		if err := recover(); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				panic(fmt.Errorf("rollback error: %v, inner error: %v", rbErr, err))
			}
			panic(err)
		}
	}()

	if err != nil {
		return fmt.Errorf("cannot begin transaction: %v", err)
	}

	if err = f(injectTx(ctx, tx)); err != nil {
		rbErr := tx.Rollback()
		return fmt.Errorf("rollback error: %v, inner error: %v", rbErr, err)
	}

	return tx.Commit()
}
