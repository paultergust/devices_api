package db

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolWithRetry(ctx context.Context) (*pgxpool.Pool, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is not set")
    }

    var pool *pgxpool.Pool
    var err error

    for i := 0; i < 10; i++ {
        pool, err = pgxpool.New(ctx, dbURL)
        if err == nil {
            err = pool.Ping(ctx)
            if err == nil {
                return pool, nil
            }
        }

        time.Sleep(2 * time.Second)
    }

    return nil, fmt.Errorf("could not connect to database: %w", err)
}
