package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	user, statusCode := service.Login(w, r)

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := service.RemoveSession(w, r)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// "/api/me" is a common RESTful convention that returns the currently authenticated user's info
func HandleMe(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := repository.GetUserById(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("00")
		return
	}

	errMsg, statusCode := service.RegisterUser(r)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		http.Error(w, errMsg, statusCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func HandleUpdateMe(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("updating profile")

	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usr, errMsg, errStatus := service.UpdateUserProfile(userId, r)
	if errMsg != "" {
		http.Error(w, errMsg, errStatus)
		fmt.Println("Error 3", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		http.Error(w, "Users not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetGroups(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetAllGroups()
	if err != nil {
		http.Error(w, "Users not found", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usr, statusCode := service.UserById(r)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

// SearchUsers handles the user search request.
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode([]model.User{}) // Return empty array for empty query
		return
	}

	users, err := repository.SearchUsers(query)
	if err != nil {
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := repository.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// handleCreatePost adds a post to the database and returns the new one
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post, statusCode := service.CreatePost(r, userID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
