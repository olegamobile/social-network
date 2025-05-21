package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
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
    n.group_id,
    g.title AS group_title,
    n.event_id,
    e.title AS event_title,
    n.content,
    n.is_read,
	n.created_at
FROM notifications n
LEFT JOIN follow_requests fr ON n.follow_req_id = fr.id
LEFT JOIN users u ON fr.follower_id = u.id
LEFT JOIN users fu ON fr.follower_id = fu.id AND n.type = 'follow_request'
LEFT JOIN group_invitations gi ON n.group_invite_id = gi.id
LEFT JOIN users iu ON gi.inviter_id = iu.id AND n.type = 'group_invitation'
LEFT JOIN group_members gm ON n.group_id = gm.group_id AND gm.user_id = n.user_id AND n.type = 'group_join_request'
LEFT JOIN users gu ON gm.user_id = gu.id AND n.type = 'group_join_request'
LEFT JOIN groups g ON n.group_id = g.id
LEFT JOIN events e ON n.event_id = e.id
LEFT JOIN users eu ON e.creator_id = eu.id AND n.type = 'event_creation'
WHERE n.status = 'enable' AND n.user_id = ?
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
