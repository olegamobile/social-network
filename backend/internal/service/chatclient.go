package service

import (
	"backend/internal/model"
	"log"
)

func ReadPump(c *model.Client) {

	defer func() {
		log.Println("closing connection for", c.UserID)

		// Remove client from map
		model.Mu.Lock()
		delete(model.Clients, c.UserID)
		model.Mu.Unlock()

		c.Conn.Close() // Close the WebSocket connection
		close(c.Send)  // Close send channel to stop WritePump
	}()

	for {
		var msg model.WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read error:", err)
			break
		}

		msg.From = c.UserID

		// Route message
		switch msg.Type {
		case "chat_message", "notification":
			//fmt.Println("normal chat message received:", msg)
			err := SaveMessage(msg)
			if err != nil {
				log.Println("failed to save message:", err)
				continue // Don't broadcast if saving failed
			}
			model.Broadcast <- msg // Send to central dispatcher
		case "groupchat_message":
			//fmt.Println("groupchat message received:", msg)
			err := SaveGroupMessage(msg)
			if err != nil {
				log.Println("failed to save group message:", err)
				continue
			}
			model.Broadcast <- msg // Send to central dispatcher
		case "typing":
			model.Broadcast <- msg
		case "ping":
			log.Println("ping from", c.UserID)
		default:
			log.Println("unknown message type:", msg.Type)
		}
	}
}

func WritePump(c *model.Client) {
	defer c.Conn.Close()

	for msg := range c.Send {
		//fmt.Println("message at write pump:", msg)

		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
