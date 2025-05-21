package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"net/http"
)

func CreatePost(payload struct {
	Content string `json:"content"`
	Privacy string `json:"privacy_level"`
}, userID int) (model.Post, int) {

	var post model.Post
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

	return post, http.StatusOK
}
