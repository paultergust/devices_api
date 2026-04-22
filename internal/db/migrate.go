package db

import (
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string) error {
    if databaseURL == "" {
        return fmt.Errorf("DATABASE_URL is empty")
    }

    m, err := migrate.New(
        "file://migrations",
        databaseURL,
    )
    if err != nil {
        return err
    }

    err = m.Up()
    if err != nil && err.Error() != "no change" {
        return err
    }

    return nil
}
