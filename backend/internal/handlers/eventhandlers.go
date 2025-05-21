package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	// "fmt"
	"net/http"
	// "strconv"
	// "strings"
)

func HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	event.CreatorID = userID

	id, err := repository.CreateEvent(event)
	if err != nil {
		http.Error(w, "Could not create event", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"event_id": id})
}

func HandleEventResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var resp model.EventResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp.UserID = userID

	if resp.Response != "going" && resp.Response != "not_going" && resp.Response != "pending" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	err = repository.SaveEventResponse(resp.EventID, resp.UserID, resp.Response)
	if err != nil {
		http.Error(w, "Could not save response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
