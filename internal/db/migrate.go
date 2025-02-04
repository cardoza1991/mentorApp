package db

import (
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDB handles database migrations
func MigrateDB(dbURL string) error {
    m, err := migrate.New(
        "file://migrations",
        dbURL,
    )
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %v", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to apply migrations: %v", err)
    }

    log.Println("Database migrations applied successfully")
    return nil
}

// RollbackDB rolls back the last migration
func RollbackDB(dbURL string) error {
    m, err := migrate.New(
        "file://migrations",
        dbURL,
    )
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %v", err)
    }
    defer m.Close()

    if err := m.Steps(-1); err != nil {
        return fmt.Errorf("failed to rollback migration: %v", err)
    }

    log.Println("Database rollback successful")
    return nil
}
