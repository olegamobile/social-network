package database

import (
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

/* func runSqlFile(db *sql.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Split statements on semicolon
	//queries := strings.SplitSeq(string(sqlBytes), ";")
	queries := strings.Split(string(sqlBytes), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing query: %q\n%v", query, err)
		}
	}
	return nil
} */

/* func NewDatabase(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	//err = runInitSQL(DB, "migrations/init.sql")
	err = runSqlFile(DB, "migrations/db_preliminary_structure.sql")
	if err != nil {
		return err
	}

	err = runSqlFile(DB, "migrations/insert_data.sql")
	if err != nil {
		return err
	}

	return nil
} */

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
	if err = applyMigrations(DB, path); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database connection established and migrations applied successfully")
	return nil
}

// applyMigrations runs all database migrations from the specified directory
func applyMigrations(db *sql.DB, dbPath string) error {
	// Create a new migration driver instance
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Use the file:// driver to read migration files from the filesystem
	// Point to the migrations directory
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "sqlite3", driver)
	if err != nil {

		//return fmt.Errorf("failed to create migration instance: %w", err)
		return fmt.Errorf("failed to create migration instance")
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
