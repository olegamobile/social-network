package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ValidateSession extracts the session cookie, validates it, and returns the user_id
func ValidateSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, errors.New("missing session cookie")
	}

	userID, expiresAt, err := repository.GetUserIdAndExpirationBySessionId(cookie)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid session")
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func CreateSession(user model.User, w http.ResponseWriter) error {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	err := repository.InsertSession(sessionID, user, expiresAt)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func RemoveSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		err = repository.DeleteSessionById(cookie.Value) // proceed to wipe cookie regardless
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // expired
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return err
}
