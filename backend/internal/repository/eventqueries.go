package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	// "database/sql"
	// "fmt"
	// "net/http"
)

func CreateEvent(event model.Event) (int, error) {
	query := `INSERT INTO events (group_id, creator_id, title, description, event_datetime)
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate)
	if err != nil {
		return 0, err
	}
	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Loop throung the group members, crete event responses and send notifications
	members, err := GetGroupMembersByGroupId(int(eventID))
	if err != nil {
		return 0, err
	}
	for _, member := range members {
		err = SaveEventResponse(int(eventID), member.ID, "pending")
		if err != nil {
			return 0, err
		}
		// send notifications
		err = InsertNotification(event.CreatorID, member.ID, "event_creation", *event.ID)
		if err != nil {
			return 0, err
		}
	}

	return int(eventID), err
}

func SaveEventResponse(eventID, userID int, response string) error {
	query := `
INSERT INTO event_responses (event_id, user_id, response)
VALUES (?, ?, ?)
ON CONFLICT(event_id, user_id) DO UPDATE SET response = ?
	`
	_, err := database.DB.Exec(query, eventID, userID, response)
	return err
}
