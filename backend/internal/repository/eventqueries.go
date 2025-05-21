package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	// "database/sql"
	// "fmt"
	// "net/http"
)

func CreateEvent(event model.Event) (int, error) {
	query := `INSERT INTO events (group_id, creator_id, title, description, event_datetime, created_at, status)
	          VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?)`
	result, err := database.DB.Exec(query, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate, event.Status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func SaveEventResponse(eventID, userID int, response string) error {
	query := `
	INSERT INTO event_responses (event_id, user_id, response)
	VALUES (?, ?, ?)
	ON CONFLICT(event_id, user_id) DO UPDATE SET response = excluded.response
	`
	_, err := database.DB.Exec(query, eventID, userID, response)
	return err
}