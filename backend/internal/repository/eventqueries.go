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
	members, err := GetGroupMembersByGroupId(int(event.GroupID))
	if err != nil {
		return 0, err
	}
	for _, member := range members {
		err = SaveEventResponse(int(eventID), member.ID, "pending")
		if err != nil {
			return 0, err
		}
		// send notifications
		_, err = InsertNotification(event.CreatorID, member.ID, "event_creation", *event.ID)
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

func GetEventByID(eventID int) (model.Event, error) {
	var e model.Event
	query := `
		SELECT e.id, e.group_id, g.title, e.creator_id, e.title, e.description, e.event_datetime
		FROM events e
		LEFT JOIN groups g ON e.group_id = g.id
		WHERE e.id = ? AND e.status = 'enable'
	`
	err := database.DB.QueryRow(query, eventID).Scan(&e.ID, &e.GroupID, &e.Group, &e.CreatorID, &e.Title, &e.Description, &e.EventDate)
	return e, err
}

func CheckUserGroupMembership(userID, groupID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM group_members
			WHERE user_id = ? AND group_id = ? AND status = 'enable' AND approval_status = 'accepted'
		)
	`
	err := database.DB.QueryRow(query, userID, groupID).Scan(&exists)
	return exists, err
}

func GetEventsByUser(userID int) ([]model.Event, error) {
	query := `
	SELECT DISTINCT e.id, g.title, e.group_id, e.creator_id, e.title, e.description, e.event_datetime
	FROM events e
	JOIN groups g ON e.group_id = g.id
	LEFT JOIN event_responses er ON er.event_id = e.id
	WHERE e.creator_id = ? OR er.user_id = ?
	ORDER BY e.event_datetime DESC
	`
	rows, err := database.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		err := rows.Scan(&e.ID, &e.Group, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventDate)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func GetEventsByGroup(groupID, userID int) ([]model.Event, error) {
	query := `
	SELECT DISTINCT e.id, g.title, e.group_id, e.creator_id, e.title, e.description, e.event_datetime
	FROM events e
	JOIN groups g ON e.group_id = g.id
	LEFT JOIN event_responses er ON er.event_id = e.id
	WHERE e.group_id = ?
	  AND e.status = 'enable'
	  AND (e.creator_id = ? OR er.user_id = ?)
	ORDER BY e.event_datetime DESC
	`

	rows, err := database.DB.Query(query, groupID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		err := rows.Scan(&e.ID, &e.Group, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventDate)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
