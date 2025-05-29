package service

import (
	"backend/internal/model"
	"encoding/json"
	"log"
	"strconv"
)

// SendNotificationDeletedWS sends a WebSocket message to notify a user that a notification has been deleted.
func SendNotificationDeletedWS(notificationID int64, recipientUserID int) error {
	content := map[string]string{"id": strconv.FormatInt(notificationID, 10)}
	jsonContent, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error marshalling notification deletion content: %v", err)
		return err
	}

	wsMsg := model.WSMessage{
		Type:    "notification_deleted",
		To:      strconv.Itoa(recipientUserID),
		From:    "system_notification_deletion",
		Content: string(jsonContent),
	}

	select {
	case model.Broadcast <- wsMsg:
		log.Printf("Notification deletion message for notification ID %d queued for user %s.", notificationID, wsMsg.To)
	default:
		log.Printf("Broadcast channel full or unavailable. Notification deletion message for notification ID %d for user %s not sent via WebSocket immediately.", notificationID, wsMsg.To)
		// Depending on requirements, you might return an error here:
		// return fmt.Errorf("broadcast channel full, message for notification %d not sent", notificationID)
	}

	return nil
}
