package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	"strconv"
)

func SaveMessage(msg model.WSMessage) error {
	db := database.DB
	var err error

	_, err = db.Exec(`
        INSERT INTO messages (sender_id, receiver_id, content)
        VALUES (?, ?, ?)
    `, msg.From, msg.To, msg.Content)

	return err
}

func GetUserChats(userId int) ([]model.Chat, error) {
	query := `
	SELECT
	  m.id,
	  m.sender_id,
	  m.receiver_id,
	  m.content,
	  m.created_at,
	  m.updated_at,
	  CASE
	    WHEN m.sender_id = ? THEN m.receiver_id
	    ELSE m.sender_id
	  END AS other_user_id,
	  u.first_name,
	  u.last_name,
	  EXISTS (
	    SELECT 1 FROM follow_requests fr
	    WHERE
	      fr.approval_status = 'accepted' AND (
	        (fr.follower_id = ? AND fr.followed_id = 
	          CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END) OR
	        (fr.follower_id = 
	          CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END AND fr.followed_id = ?)
	      )
	  ) AS is_active
	FROM messages m
	JOIN users u
	  ON u.id = CASE
	    WHEN m.sender_id = ? THEN m.receiver_id
	    ELSE m.sender_id
	  END
	WHERE
	  (m.sender_id = ? OR m.receiver_id = ?)
	  AND m.status = 'enable'
	ORDER BY m.created_at ASC;
	`

	rows, err := database.DB.Query(query, userId, userId, userId, userId, userId, userId, userId, userId)
	if err != nil {
		fmt.Println("query error at GetUserChats:", err)
		return nil, err
	}
	defer rows.Close()

	chatMap := make(map[int]*model.Chat)

	for rows.Next() {
		var (
			msg                 model.ChatMessage
			otherUserID         int
			firstName, lastName string
			updatedAt           sql.NullString
			isActive            bool
		)

		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &updatedAt, &otherUserID, &firstName, &lastName, &isActive)

		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			msg.UpdatedAt = updatedAt.String
		}

		err = database.DB.QueryRow("SELECT first_name FROM users WHERE id = ?", msg.SenderID).Scan(&msg.SenderName)
		if err != nil {
			return nil, err
		}

		if _, exists := chatMap[otherUserID]; !exists {
			chatMap[otherUserID] = &model.Chat{
				Name:     fmt.Sprintf("%s %s", firstName, lastName),
				UserID:   strconv.Itoa(otherUserID),
				Messages: []model.ChatMessage{},
				IsActive: isActive,
			}
		}
		chatMap[otherUserID].Messages = append(chatMap[otherUserID].Messages, msg)
	}

	// Convert map to slice
	chats := make([]model.Chat, 0, len(chatMap))
	for _, chat := range chatMap {
		chats = append(chats, *chat)
	}

	return chats, nil
}
