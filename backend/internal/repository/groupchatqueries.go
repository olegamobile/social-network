package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"fmt"
)

func GetGroupChat(groupID int) (model.Chat, error) {
	var chat model.Chat

	query := `
	SELECT gm.sender_id, u.first_name, gm.content, gm.created_at
	FROM group_messages gm
	JOIN users u ON gm.sender_id = u.id
	WHERE gm.group_id = ? AND gm.status = 'enable'
	`

	rows, err := database.DB.Query(query, groupID)
	if err != nil {
		fmt.Println("query error in GetGroupChat", err)
		return chat, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.ChatMessage

		err := rows.Scan(&msg.SenderID, &msg.SenderName, &msg.Content, &msg.CreatedAt)
		if err != nil {
			fmt.Println("scan error in GetGroupChat", err)
			return chat, err
		}
		chat.Messages = append(chat.Messages, msg)
	}

	chat.IsActive = true
	chat.UserID = fmt.Sprint(groupID)

	return chat, nil
}

func SaveGroupMessage(msg model.WSMessage) error {
	db := database.DB
	var err error

	_, err = db.Exec(`
        INSERT INTO group_messages (sender_id, group_id, content)
        VALUES (?, ?, ?)
    `, msg.From, msg.To, msg.Content)

	return err
}

/*
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   group_id INTEGER NOT NULL,
   sender_id INTEGER NOT NULL,
   content TEXT NOT NULL,
   created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME,
   updated_by INTEGER,
   status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable', */
