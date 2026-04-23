package db

import (
    "context"
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    ErrMissingDatabaseURL   = errors.New("database url is not set")
    ErrDatabaseUnavailable  = errors.New("database is unavailable")
)

func NewPoolWithRetry(ctx context.Context) (*pgxpool.Pool, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, ErrMissingDatabaseURL
    }

    var pool *pgxpool.Pool
    var lastErr error

    for i := 0; i < 10; i++ {
        // respect cancellation
        select {
        case <-ctx.Done():
            return nil, fmt.Errorf("context canceled while connecting to database: %w", ctx.Err())
        default:
        }

        pool, lastErr = pgxpool.New(ctx, dbURL)
        if lastErr != nil {
            lastErr = fmt.Errorf("creating pool: %w", lastErr)
        } else {
            if err := pool.Ping(ctx); err == nil {
                return pool, nil
            } else {
                lastErr = fmt.Errorf("ping database: %w", err)
                pool.Close() // avoid leaking connections on retry
            }
        }

        time.Sleep(2 * time.Second)
    }

    return nil, fmt.Errorf("%w after retries: %v", ErrDatabaseUnavailable, lastErr)
}
