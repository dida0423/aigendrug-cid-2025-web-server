package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ParsePostgresConnectionString(connectionString string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func WithTx(ctx context.Context, db *pgxpool.Pool, fn func(tx pgx.Tx) error) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("panic recovered in transaction: %v", p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(tx)
	return
}

func WithTxResult[T any](
	ctx context.Context,
	db *pgxpool.Pool,
	fn func(tx pgx.Tx) (T, error),
) (result T, err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("panic recovered in transaction: %v", p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	result, err = fn(tx)
	return
}

func TimeToPgTimestamptz(t time.Time) (pgtype.Timestamptz, error) {
	var ts pgtype.Timestamptz
	err := ts.Scan(t)
	if err != nil {
		return pgtype.Timestamptz{}, err
	}
	return ts, nil
}
