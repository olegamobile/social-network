package repository

import (
	"backend/internal/database"
	"backend/internal/model"
)

func SaveMessage(msg model.WSMessage) error {
	db := database.DB
	var err error

	_, err = db.Exec(`
        INSERT INTO messages (sender_id, receiver_id, content)
        VALUES (?, ?, ?)
    `, msg.To, msg.From, msg.Content)

	return err
}
