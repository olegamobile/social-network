package handlers

import (
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"net/http"
)

func HandleGetUserMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chats, err := repository.GetUserChats(userID)
	if err != nil {
		//fmt.Println("error getting chats at HandleGetUserMessages", err)
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
