package model

import "database/sql"

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Birthday   string `json:"birthday"`
	Password   string `json:"password"`
	About      string `json:"about_me"`
	AvatarPath string `json:"avatar_url"`
}

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileData struct {
	FirstName string
	LastName  string
	DOB       string
	Nickname  string
	About     string
	//AvatarPath   *string
	AvatarPath   sql.NullString
	DeleteAvatar bool
}

type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
