package db

import (
    "errors"
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
    ErrMigrationFailed    = errors.New("migration failed")
)

func RunMigrations(databaseURL string) error {
    if databaseURL == "" {
        return ErrMissingDatabaseURL
    }

    m, err := migrate.New(
        "file://migrations",
        databaseURL,
    )
    if err != nil {
        return fmt.Errorf("initialize migrate instance: %w", err)
    }

    // Ensure resources are cleaned up
    defer func() {
        _, _ = m.Close()
    }()

    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            return nil
        }
        return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
    }

    return nil
}
