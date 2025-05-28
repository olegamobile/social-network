package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

func GetAllNotificatons(userID int) ([]model.Notification, error) {
	rows, err := database.DB.Query(`
SELECT 
    n.id,
    n.type,
    n.user_id,
    CASE 
        WHEN n.type = 'follow_request' THEN fr.follower_id
        WHEN n.type = 'group_invitation' THEN gi.inviter_id
        WHEN n.type = 'group_join_request' THEN gm.user_id
        WHEN n.type = 'event_creation' THEN e.creator_id
        ELSE NULL
    END AS sender_id,
    CASE 
        WHEN n.type = 'follow_request' THEN (fu.first_name || ' ' || fu.last_name)
        WHEN n.type = 'group_invitation' THEN (iu.first_name || ' ' || iu.last_name)
        WHEN n.type = 'group_join_request' THEN (gu.first_name || ' ' || gu.last_name)
        WHEN n.type = 'event_creation' THEN (eu.first_name || ' ' || eu.last_name)
        ELSE NULL
    END AS sender_name,
    n.follow_req_id,
    n.group_invite_id,

	CASE 
        WHEN n.type = 'group_invitation' THEN gi.group_id
        WHEN n.type = 'group_join_request' THEN gm.group_id
        ELSE NULL
    END AS group_id,
    
    -- Group title selection based on type
    COALESCE(ggm.title, ggi.title, ge.title) AS group_title,

    n.event_id,
    e.title AS event_title,
    n.content,
    n.is_read,
	
	CASE 
        WHEN n.updated_at IS NULL THEN n.created_at
        ELSE n.updated_at
    END AS notification_time

FROM notifications n
LEFT JOIN follow_requests fr ON n.follow_req_id = fr.id
LEFT JOIN users u ON fr.follower_id = u.id
LEFT JOIN users fu ON fr.follower_id = fu.id AND n.type = 'follow_request'
LEFT JOIN group_invitations gi ON n.group_invite_id = gi.id
LEFT JOIN users iu ON gi.inviter_id = iu.id AND n.type = 'group_invitation'
LEFT JOIN groups ggi ON gi.group_id = ggi.id AND n.type = 'group_invitation'
LEFT JOIN group_members gm ON n.group_members_id = gm.id AND n.type = 'group_join_request'
LEFT JOIN users gu ON gm.user_id = gu.id AND n.type = 'group_join_request'
LEFT JOIN groups ggm ON gm.group_id = ggm.id AND n.type = 'group_join_request'
LEFT JOIN events e ON n.event_id = e.id
LEFT JOIN users eu ON e.creator_id = eu.id AND n.type = 'event_creation'
LEFT JOIN groups ge ON e.group_id = ge.id AND n.type = 'event_creation'
WHERE n.status = 'enable' AND n.user_id = ?
ORDER BY notification_time DESC
	`, userID)
	if err != nil {
		fmt.Println("query error at GetAllNotificatons", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.UserID,
			&n.SenderID,
			&n.SenderName,
			&n.FollowReqID,
			&n.GroupInviteID,
			&n.GroupID,
			&n.GroupTitle,
			&n.EventID,
			&n.EventTitle,
			&n.Content,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			fmt.Println("scan error at GetAllNotificatons", err)
			return notifications, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func GetNotification(userID int, notificationID int) (model.Notification, error) {
	var n model.Notification
	err := database.DB.QueryRow(`
        SELECT id, type, user_id, follow_req_id, group_invite_id, event_id, content
        FROM notifications
        WHERE id = ? AND user_id = ? AND status = 'enable'
    `, notificationID, userID).Scan(&n.ID, &n.Type, &n.UserID, &n.FollowReqID, &n.GroupInviteID, &n.EventID, &n.Content)
	if err != nil {
		return n, err
	}
	return n, nil
}

func MarkNotificationAsRead(userID int, notificationID int) error {
	_, err := database.DB.Exec(`
        UPDATE notifications
        SET is_read = 1, updated_at = CURRENT_TIMESTAMP, updated_by = ?
        WHERE id = ? AND user_id = ? AND status = 'enable'
    `, userID, notificationID, userID)
	return err
}

func CheckFollowRequestStatus(followRequestID int) (bool, error) {
	status := ""
	err := database.DB.QueryRow(`
        SELECT approval_status
        FROM follow_requests
        WHERE id = ?
    `, followRequestID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows { // no row found means not pending
			return false, nil
		}
		return false, err
	}
	return status == "pending", nil
}

func CheckInvitationStatus(groupInvitationID int) (bool, error) {
	status := ""
	err := database.DB.QueryRow(`
        SELECT approval_status
        FROM group_invitations
        WHERE id = ?
    `, groupInvitationID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows { // no row found means not pending
			return false, nil
		}
		return false, err
	}
	return status == "pending", nil
}

func CheckJoinRequestStatus(userID, groupID int) (bool, error) {
	status := ""
	err := database.DB.QueryRow(`
        SELECT approval_status
        FROM group_members
        WHERE user_id = ? AND group_id = ?
    `, userID, groupID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows { // no row found means not pending
			return false, nil
		}
		return false, err
	}
	return status == "pending", nil
}

func CheckEventInvitationStatus(userID, eventID int) (bool, error) {
	status := ""
	err := database.DB.QueryRow(`
        SELECT response
        FROM event_responses
        WHERE user_id = ? AND event_id = ? 
    `, userID, eventID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows { // no row found means not pending
			return false, nil
		}
		return false, err
	}
	return status == "pending", nil
}

func GetNewNotificatonsCount(userID int) (int, error) {
	count := 0
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = ? AND is_read = 0 AND status = 'enable'
	`, userID).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func GetJoinRequests(groupId int) ([]model.Notification, error) {
	query := `
	SELECT n.id, n.type, n.user_id, gm.user_id, u.first_name || ' ' || u.last_name, gm.group_id, g.title, n.content, n.is_read, n.created_at
	FROM notifications n
	JOIN group_members gm ON n.group_members_id = gm.id
	JOIN groups g ON gm.group_id = g.id
	JOIN users u ON gm.user_id = u.id
	WHERE gm.group_id = ? AND n.type = 'group_join_request' AND n.status = 'enable';
	`

	rows, err := database.DB.Query(query, groupId)
	if err != nil {
		fmt.Println("query error at GetJoinRequests", err)
		return nil, err
	}

	var notices []model.Notification
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.UserID,
			&n.SenderID,
			&n.SenderName,
			&n.GroupID,
			&n.GroupTitle,
			&n.Content,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			fmt.Println("scan error at GetAllNotificatons", err)
			return notices, err
		}
		notices = append(notices, n)
	}

	return notices, nil
}

// InsertNotification handles creating or updating a notification and then broadcasts it.
func InsertNotification(userFromID, userToID int, notifType string, refID int) (int64, error) {
	var notificationID int64
	var dbErr error

	// Note: userFromID is the initiator of the action causing the notification.
	// userToID is the recipient of the notification.

	now := time.Now().Format("2006-01-02 15:04:05") // For created_at and updated_at

	// Base part of the query for insertion
	// The specific reference ID column (follow_req_id, group_invite_id, etc.) is set based on notifType
	var insertColumnName string
	switch notifType {
	case "follow_request":
		insertColumnName = "follow_req_id"
	case "group_invitation":
		insertColumnName = "group_invite_id"
	case "group_join_request":
		insertColumnName = "group_members_id" // Corresponds to group_members.id
	case "event_creation":
		insertColumnName = "event_id"
	default:
		return 0, fmt.Errorf("invalid notification type: %s", notifType)
	}

	// Construct the full query
	// The ON CONFLICT targets the (user_id, ref_type, ref_id) unique constraint.
	// ref_type and ref_id are generated columns based on type and the specific ID columns.
	query := fmt.Sprintf(`
		INSERT INTO notifications (user_id, type, %s, created_at, is_read, status)
		VALUES (?, ?, ?, ?, FALSE, 'enable')
		ON CONFLICT(user_id, ref_type, ref_id) 
		DO UPDATE SET 
			updated_at = ?,
			updated_by = ?, 
			status = 'enable',
			is_read = FALSE
		RETURNING id`, insertColumnName)

	// Execute the query
	// For INSERT: userToID, notifType, refID, now
	// For UPDATE: now, userFromID (as updated_by)
	dbErr = database.DB.QueryRow(query, userToID, notifType, refID, now, now, userFromID).Scan(&notificationID)

	if dbErr != nil {
		log.Printf("Error inserting/updating notification: %v. Query: %s, Args: userToID=%d, notifType=%s, refID=%d, now=%s, updated_by_userFromID=%d", dbErr, query, userToID, notifType, refID, now, userFromID)
		return 0, dbErr
	}

	// If DB operation was successful, proceed to broadcast
	if notificationID > 0 {
		fullNotification, broadcastErr := GetNotificationForBroadcasting(notificationID)
		if broadcastErr != nil {
			log.Printf("Error fetching notification %d for broadcasting: %v", notificationID, broadcastErr)
			// Do not return this error, primary operation (DB insert) succeeded
		} else if fullNotification == nil {
			log.Printf("Notification %d not found for broadcasting (was nil)", notificationID)
		} else {
			jsonData, marshalErr := json.Marshal(fullNotification)
			if marshalErr != nil {
				log.Printf("Error marshalling notification %d for broadcasting: %v", notificationID, marshalErr)
			} else {
				wsMsg := model.WSMessage{
					Type:    "new_notification",
					To:      strconv.Itoa(userToID),
					From:    "system_notification",
					Content: string(jsonData),
				}
				
				// Non-blocking send to broadcast channel
				if model.Broadcast != nil {
					select {
					case model.Broadcast <- wsMsg:
						log.Printf("Notification %d queued for user %s.", notificationID, wsMsg.To)
					default:
						log.Printf("Broadcast channel full or unavailable. Notification %d for user %s not sent via WebSocket immediately.", notificationID, wsMsg.To)
					}
				} else {
					log.Println("model.Broadcast channel is nil. WebSocket message not sent.")
				}
			}
		}
	}

	return notificationID, nil // Return original DB error (nil if successful)
}

func GetNotificationForBroadcasting(notificationID int64) (*model.Notification, error) {
	var n model.Notification
	query := `
		SELECT
            n.id, n.type, n.user_id,
            CASE
                WHEN n.type = 'follow_request' THEN fr.follower_id
                WHEN n.type = 'group_invitation' THEN gi.inviter_id
                WHEN n.type = 'group_join_request' THEN gm.user_id
                WHEN n.type = 'event_creation' THEN e.creator_id
                ELSE NULL
            END AS sender_id,
            CASE
                WHEN n.type = 'follow_request' THEN (fu.first_name || ' ' || fu.last_name)
                WHEN n.type = 'group_invitation' THEN (iu.first_name || ' ' || iu.last_name)
                WHEN n.type = 'group_join_request' THEN (gu.first_name || ' ' || gu.last_name)
                WHEN n.type = 'event_creation' THEN (eu.first_name || ' ' || eu.last_name)
                ELSE NULL
            END AS sender_name,
            n.follow_req_id, n.group_invite_id,
            CASE
                WHEN n.type = 'group_invitation' THEN gi.group_id
                WHEN n.type = 'group_join_request' THEN gm.group_id
                ELSE NULL
            END AS group_id,
            COALESCE(ggm.title, ggi.title, ge.title) AS group_title,
            n.event_id, e.title AS event_title,
            n.content, n.is_read,
            strftime('%Y-%m-%d %H:%M:%S', COALESCE(n.updated_at, n.created_at)) AS notification_time
        FROM notifications n
        LEFT JOIN follow_requests fr ON n.follow_req_id = fr.id
        LEFT JOIN users fu ON fr.follower_id = fu.id AND n.type = 'follow_request'
        LEFT JOIN group_invitations gi ON n.group_invite_id = gi.id
        LEFT JOIN users iu ON gi.inviter_id = iu.id AND n.type = 'group_invitation'
        LEFT JOIN groups ggi ON gi.group_id = ggi.id AND n.type = 'group_invitation'
        LEFT JOIN group_members gm ON n.group_members_id = gm.id AND n.type = 'group_join_request'
        LEFT JOIN users gu ON gm.user_id = gu.id AND n.type = 'group_join_request'
        LEFT JOIN groups ggm ON gm.group_id = ggm.id AND n.type = 'group_join_request'
        LEFT JOIN events e ON n.event_id = e.id
        LEFT JOIN users eu ON e.creator_id = eu.id AND n.type = 'event_creation'
        LEFT JOIN groups ge ON e.group_id = ge.id AND n.type = 'event_creation'
        WHERE n.id = ? AND n.status = 'enable'
	`
	err := database.DB.QueryRow(query, notificationID).Scan(
		&n.ID,
		&n.Type,
		&n.UserID,
		&n.SenderID,
		&n.SenderName,
		&n.FollowReqID,
		&n.GroupInviteID,
		&n.GroupID,
		&n.GroupTitle,
		&n.EventID,
		&n.EventTitle,
		&n.Content,
		&n.IsRead,
		&n.CreatedAt, // This corresponds to notification_time from the query
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Specific error for not found
		}
		return nil, fmt.Errorf("error scanning notification for broadcasting: %w", err)
	}

	return &n, nil
}
