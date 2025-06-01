package service

import (
	"backend/internal/model"
	"log"
)

func StartBroadcastListener() {
	//go func() {	// starts as goroutine already in main.go
	for {
		msg := <-model.Broadcast
		model.Mu.Lock()
		for _, client := range model.Clients {
			if msg.To == "" || client.UserID == msg.To {
				select {
				case client.Send <- msg:
				default:
					log.Printf("Send buffer full for user %s, skipping message", client.UserID)
				}
			}
		}
		model.Mu.Unlock()
	}
	//}()
}
