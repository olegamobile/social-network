package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreatePost(r *http.Request, userID int) (model.Post, int) {
	var post model.Post

	var payload struct {
		Content string `json:"content"`
		Privacy string `json:"privacy_level"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.Content == "" {
		fmt.Println("error in CreatePost:", err)
		return post, http.StatusBadRequest
	}

	id, createdAt, err := repository.InsertPost(userID, payload.Content, payload.Privacy)
	if err != nil {
		return post, http.StatusInternalServerError
	}

	usr, err := repository.GetUserById(userID, true)
	if err != nil {
		return post, http.StatusNotFound
	}

	post.ID = id
	post.UserID = userID
	post.Username = usr.FirstName + " " + usr.LastName
	post.AvatarPath = usr.AvatarPath
	post.Content = payload.Content
	post.CreatedAt = createdAt
	post.PostType = "regular"
	return post, http.StatusOK
}
