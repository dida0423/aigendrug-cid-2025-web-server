package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(c context.Context) (*pgxpool.Pool, error) {
	config, err := ParsePostgresConnectionString(os.Getenv("DB_CONNECTION_STRING"))

	if err != nil {
		return nil, err
	}

	AutoMigrateFromConnectionString(c, os.Getenv("DB_CONNECTION_STRING"), config)

	dbpool, err := pgxpool.NewWithConfig(c, config)

	if err != nil {
		print("(pool.go) Failed to connect to database: ", err)
		return nil, err
	}

	return dbpool, nil
}

// func NewPostgresPoolFromConnectionString(c context.Context, connectionString string) (*pgxpool.Pool, error) {
// 	config, err := ParsePostgresConnectionString(connectionString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return NewPostgresPool(c, config)
// }
