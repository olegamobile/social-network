package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandleFollowing(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/following/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	followStatus, statusCode := service.GetFollowSatatus(userId, targetId)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followStatus)
}

func HandleFollowAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleFollowAction")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("ValidateSession error at HandleFollowAction:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.FollowRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("json error at HandleFollowAction:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var statusCode int
	switch req.Action {
	case "request":
		statusCode = service.FollowRequest(userID, req.TargetID)
	case "follow":
		statusCode = service.Follow(userID, req.TargetID)
	case "unfollow":
		statusCode = service.Unfollow(userID, req.TargetID)
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
	}

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleFollowAction:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/followers/"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	users, err := repository.GetFollowersByUserID(userId)
	if err != nil {
		http.Error(w, "Error getting followers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetFollowedUsers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/followed/"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	users, err := repository.GetFollowedUsersByUserID(userId)
	if err != nil {
		http.Error(w, "Error getting followed users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetSentFollowRequests(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := repository.GetFollowRequestsSentByUser(userId)
	if err != nil {
		http.Error(w, "Error getting sent follow requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetReceivedFollowRequests(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := repository.GetFollowRequestsReceivedByUser(userId)
	if err != nil {
		http.Error(w, "Error getting received follow requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
