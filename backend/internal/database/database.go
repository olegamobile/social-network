package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

var DB *sql.DB

func runInitSQL(db *sql.DB, filepath string) error {
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
}

func NewDatabase(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	err = runInitSQL(DB, "migrations/init.sql")
	if err != nil {
		return err
	}

	return nil
}
