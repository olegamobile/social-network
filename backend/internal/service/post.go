package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"encoding/json"
	"net/http"
)

func CreatePost(r *http.Request, userID int) (model.Post, int) {
	var post model.Post

	var payload struct {
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.Content == "" {
		return post, http.StatusBadRequest
	}

	id, createdAt, err := repository.InsertPost(userID, payload.Content)
	if err != nil {
		return post, http.StatusInternalServerError
	}

	usr, err := repository.GetUserById(userID)
	if err != nil {
		return post, http.StatusNotFound
	}

	post.ID = id
	post.UserID = userID
	post.Username = usr.Username
	post.Content = payload.Content
	post.CreatedAt = createdAt

	return post, http.StatusOK
}
