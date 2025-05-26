package handlers

import (
	"backend/config"
	//"backend/internal/repository"
	"backend/internal/service"
	//"encoding/json"
	"backend/internal/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == config.FrontendURL // explicitly allow frontend origin
	},
}

func HandleWSConnections(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &model.Client{
		UserID: fmt.Sprint(userID),
		Conn:   conn,
		Send:   make(chan model.WSMessage, 10), // buffered send
	}

	model.Mu.Lock()
	model.Clients[client.UserID] = client
	model.Mu.Unlock()

	go service.ReadPump(client)
	service.WritePump(client)

	model.Mu.Lock()
	delete(model.Clients, client.UserID)
	model.Mu.Unlock()
}
