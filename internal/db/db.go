package db

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() (*pgxpool.Pool, error) {
    dbURL := os.Getenv("DATABASE_URL")

    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect db: %w", err)
    }

    return pool, nil
}
