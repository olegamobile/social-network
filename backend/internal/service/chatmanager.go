package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"log"
	"strconv"
)

func sendToClient(c *model.Client, msg model.WSMessage) {
	model.Mu.Lock()
	select {
	case c.Send <- msg:
	default:
		log.Printf("Send buffer full for user %s, skipping message", c.UserID)
	}
	model.Mu.Unlock()
}

func StartBroadcastListener() {
	//go func() {	// starts as goroutine already in main.go
	for {
		msg := <-model.Broadcast

		// try to send chat message to one person
		if msg.Type == "chat_message" && msg.To != "" {
			if client, ok := model.Clients[msg.To]; ok { // matching connection found
				sendToClient(client, msg)
			}
			continue
		}

		// send group chat message to all connected group members
		if msg.Type == "groupchat_message" && msg.To != "" {
			groupId, err := strconv.Atoi(msg.To)
			if err != nil {
				continue
			}
			members, err := repository.GetGroupMembersByGroupId(groupId)
			if err != nil {
				continue
			}
			for _, member := range members {
				if client, ok := model.Clients[fmt.Sprint(member.ID)]; ok {
					sendToClient(client, msg)
				}
			}
			continue
		}

		// if recipient still specified, assume it's user id
		if msg.To != "" {
			if client, ok := model.Clients[msg.To]; ok {
				sendToClient(client, msg)
			}
			continue
		}

		// send to all if no recipient
		for _, client := range model.Clients {
			sendToClient(client, msg)
		}
	}
	//}()
}
