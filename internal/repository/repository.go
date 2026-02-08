package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

type Transactor struct {
	DB *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		DB: db,
	}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(ctxTx context.Context) error) error {
	tx, err := t.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed begin transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %v, original error: %w", rbErr, err)
		}

		return err
	}

	return tx.Commit()
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func executorFromContext(ctx context.Context, db *sql.DB) dbExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok && tx != nil {
		return tx
	}

	return db
}
