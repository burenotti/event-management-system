package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DbContext interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
	sqlx.PreparerContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
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

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) WithContext(ctx context.Context) DbContext {
	if tx := extractTx(ctx); tx != nil {
		return tx
	} else {
		return d.db
	}
}

type AtomicFunc func(ctx context.Context) error

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
