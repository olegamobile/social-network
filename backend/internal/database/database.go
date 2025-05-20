package database

import (
	"backend/config"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// NewDatabase initializes a new SQLite database connection and runs migrations
func NewDatabase(path string) error {
	var err error

	// Open database connection
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err = applyMigrations(DB); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database connection established and migrations applied successfully")
	return nil
}

// applyMigrations runs all database migrations from the specified directory
func applyMigrations(db *sql.DB) error {
	// Create a new migration driver instance
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Use the file:// driver to read migration files from the filesystem
	m, err := migrate.NewWithDatabaseInstance("file://"+config.MigrationsPath, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Apply all up migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
