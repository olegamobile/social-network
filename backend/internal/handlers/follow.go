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

	var statusCode, frID int
	switch req.Action {
	case "request":
		frID, statusCode = service.FollowRequest(userID, req.TargetID)
		if statusCode == http.StatusOK {
			err = repository.InsertNotification(userID, req.TargetID, "follow_request", frID)
			if err != nil {
				http.Error(w, "error inserting notification in HandleGroupMembership", http.StatusBadRequest)
				return
			}
		}
	case "follow":
		statusCode = service.Follow(userID, req.TargetID)
	case "unfollow":
		statusCode = repository.RemoveFollow(userID, req.TargetID)
	case "cancel":
		statusCode = repository.RemoveFollowRequestNotification(userID, req.TargetID) // this first, other is hard delete
		if statusCode == http.StatusOK {
			statusCode = repository.RemoveFollow(userID, req.TargetID)
		}
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
		return
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

	if userId == 0 {
		userId, err = service.ValidateSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	users, err := repository.GetFollowersByUserID(userId)
	if err != nil {
		http.Error(w, "Error getting followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	if userId == 0 {
		userId, err = service.ValidateSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	users, err := repository.GetFollowedUsersByUserID(userId)
	if err != nil {
		http.Error(w, "Error getting followed users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func HandleFollowRequestApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed at HandleFollowRequestApprove")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("ValidateSession error at HandleFollowRequestApprove:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/follow/requests/"), "/")
	if len(data) != 2 {
		http.Error(w, "Invalid request action syntax", http.StatusBadRequest)
		return
	}

	action := data[1]
	requestID, err := strconv.Atoi(data[0])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	var statusCode int
	switch action {
	case "accept":
		statusCode = repository.AcceptFollowRequest(userID, requestID)
	case "decline":
		statusCode = repository.AcceptFollowRequest(userID, requestID)
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
	}

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleFollowRequestApprove:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
