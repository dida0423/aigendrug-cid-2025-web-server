package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/init.sql
var initial_sql string

func AutoMigrateFromConnectionString(ctx context.Context, connectionString string, config *pgxpool.Config) (bool, error) {
	dbName := config.ConnConfig.Database
	config.ConnConfig.Database = "postgres"

	dbAdmin, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return false, err
	}

	defer dbAdmin.Close()

	_, err = dbAdmin.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		println("Database already exists, skipping creation")
	}

	config.ConnConfig.Database = dbName
	dbTarget, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return false, err
	}
	defer dbTarget.Close()

	_, err = dbTarget.Exec(ctx, initial_sql)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code != "42P07" {
				return false, err
			}
		}
	}
	println("Database initialized successfully")

	return true, nil
}
