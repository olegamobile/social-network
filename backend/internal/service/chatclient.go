package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"log"
)

func ReadPump(c *model.Client) {
	//defer c.Conn.Close()

	defer func() {
		log.Println("closing connection for", c.UserID)

		// Remove client from map
		model.Mu.Lock()
		delete(model.Clients, c.UserID)
		model.Mu.Unlock()

		// Close the WebSocket connection
		c.Conn.Close()

		// Close send channel to stop WritePump
		close(c.Send)
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
			err := repository.SaveMessage(msg)
			if err != nil {
				log.Println("failed to save message:", err)
				continue // Don't broadcast if saving failed
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

		fmt.Println("message at write pump:", msg)

		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
