package transactionmanager

import (
	"context"
	"database/sql"
	"log/slog"
)

type txKey struct{}

type TransactionManager[T any] interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error)
}

type TransactionManagerImpl[T any] struct {
	db *sql.DB
}

func NewTransactionManager[T any](db *sql.DB) *TransactionManagerImpl[T] {
	return &TransactionManagerImpl[T]{db: db}
}

func (tm *TransactionManagerImpl[T]) InTransaction(ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T
	var err error
	_, ok := ctx.Value(txKey{}).(*sql.Tx)
	if ok {
		result, fnErr := fn(ctx)
		if fnErr != nil {
			slog.Error("Error executing fn: ", "error", fnErr)
			return zero, fnErr
		}
		return result, nil
	}
	tx, err := tm.db.BeginTx(ctx, nil)
	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err != nil {
		slog.Error("Error starting transaction: ", "error", err)
		return zero, err
	}
	defer func() {
		if err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				slog.Error("Error rolling back transaction: ", "error", rollErr)
			}
		}
	}()
	result, err := fn(txCtx)
	if err != nil {
		slog.Error("Error executing fn: ", "error", err)
		return zero, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		_ = tx.Rollback()
		slog.Error("Error committing transaction: ", "error", commitErr)
		return zero, commitErr
	}

	return result, nil
}

func GetTxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx := ctx.Value(txKey{})
	if tx == nil {
		return nil, false
	}

	return tx.(*sql.Tx), true
}
