package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var Pool *pgxpool.Pool

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/todo?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse config: %v\n", err))
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Sprintf("Unable to create connection pool: %v\n", err))
	}

	// Проверка подключения
	if err := Pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("Unable to ping database: %v\n", err))
	}
}