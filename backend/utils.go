package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// validateSession extracts the session cookie, validates it, and returns the user_id
func validateSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, errors.New("missing session cookie")
	}

	var userID int
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM sessions WHERE id = ?`
	err = db.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		fmt.Println("Query error in validateSession", err)
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid session")
		}
		return 0, err // some other DB error
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}
