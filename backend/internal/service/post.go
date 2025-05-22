package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func CreatePost(content, privacyLvl string, imagePath *string, userID int) (model.Post, int) {

	var post model.Post
	id, createdAt, err := repository.InsertPost(userID, content, privacyLvl, imagePath)
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
	post.ImagePath = imagePath
	post.Content = content
	post.CreatedAt = createdAt

	return post, http.StatusOK
}

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {

	ext := filepath.Ext(header.Filename)

	if ext == "" || !utils.IsAllowedImageExtension(ext) {
		fmt.Println("bad extension:", ext)
		return "", fmt.Errorf("illegal extension")
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	dst := filepath.Join("uploads/posts", filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return "/" + dst, nil
}
