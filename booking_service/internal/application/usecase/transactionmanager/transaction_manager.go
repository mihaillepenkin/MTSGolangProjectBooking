package transactionmanager

import (
	"context"
	"database/sql"
	"log/slog"
)

type TxKey struct{}

type TransactionManager[T any] interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error)
}

type TransactionManagerImpl[T any] struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewTransactionManager[T any](db *sql.DB) *TransactionManagerImpl[T] {
	return &TransactionManagerImpl[T]{db: db, logger: slog.Default().With("component", "transaction_manager")}
}

func (tm *TransactionManagerImpl[T]) InTransaction(ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T
	var err error
	_, ok := ctx.Value(TxKey{}).(*sql.Tx)
	if ok {
		result, fnErr := fn(ctx)
		if fnErr != nil {
			tm.logger.Error("Error executing fn: ", "error", fnErr)
			return zero, fnErr
		}
		return result, nil
	}
	tx, err := tm.db.BeginTx(ctx, nil)
	txCtx := context.WithValue(ctx, TxKey{}, tx)

	if err != nil {
		tm.logger.Error("Error starting transaction: ", "error", err)
		return zero, err
	}
	defer func() {
		if err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				tm.logger.Error("Error rolling back transaction: ", "error", rollErr)
			}
		}
	}()
	result, err := fn(txCtx)
	if err != nil {
		tm.logger.Error("Error executing fn: ", "error", err)
		return zero, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		_ = tx.Rollback()
		tm.logger.Error("Error committing transaction: ", "error", commitErr)
		return zero, commitErr
	}

	return result, nil
}

func GetTxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx := ctx.Value(TxKey{})
	if tx == nil {
		return nil, false
	}

	return tx.(*sql.Tx), true
}
