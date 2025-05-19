package handlers

import (
	"backend/config"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == config.FrontendURL // explicitly allow frontend origin
		},
	}
	Clients = make(map[string]*websocket.Conn)
	Mu      sync.Mutex
)

// Handle WebSocket connections
func HandleWSConnections(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("error validating user at HandleWSConnections", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"success": false})
		return
	}

	user, err := repository.GetUserById(userID, true)
	if err != nil {
		fmt.Println("Error getting user in ws connections:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"success": false})
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"success": false})
		return
	}
	defer conn.Close()

	userKey := fmt.Sprint(user.ID)
	Mu.Lock()
	Clients[userKey] = conn
	Mu.Unlock()

	// Send welcome message to client
	err = conn.WriteJSON(map[string]string{
		"type":    "welcome",
		"user_id": userKey,
	})
	if err != nil {
		log.Println("error sending welcome message:", err)
		return
	}

	// Start reading messages
	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			// Client disconnected or sent invalid JSON
			break
		}

		switch msg["type"] {
		case "typing":
			if toConn, ok := Clients[msg["to"]]; ok {
				toConn.WriteJSON(map[string]string{
					"msgType":  "typing",
					"userFrom": msg["from"],
				})
			}
		case "ping":
			fmt.Println("Connection", userKey, "pinged")
		case "chat_message":
			fmt.Println("Got chat message:", msg)
			// TODO: store message into database and (if possible?) send to other user.
		case "notification":
			fmt.Println("Got chat message:", msg)
			// TODO: store message into database and (if possible?) send to other user.
		default:
			log.Println("unknown message type or empty message:", msg)
		}
	}

	Mu.Lock()
	fmt.Println("deleting websocket client")
	delete(Clients, userKey)
	Mu.Unlock()
}
