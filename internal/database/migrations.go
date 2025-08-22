// Package database provides database connection management for the Mockbuster API.
package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

// RunMigrations runs database migrations using Goose.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	goose.SetBaseFS(nil)

	// Run migrations
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

// GetMigrationStatus returns the current migration status.
func GetMigrationStatus(db *sql.DB) error {
	goose.SetBaseFS(nil)

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Get current version
	current, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	slog.Info("Current database migration version", "version", current)
	return nil
}
