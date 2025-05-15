package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Birthday   string     `json:"birthday"`
	Password   string     `json:"password"`
	About      string     `json:"about_me"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  *int       `json:"updated_by"`
	Status     string     `json:"status"`
	AvatarPath string     `json:"avatar_url"`
}

type Post struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	AvatarPath string `json:"avatar_url"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
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

type Comment struct {
	ID               int        `json:"id"`
	PostId           int        `json:"post_id"`
	CommentId        int        `json:"comment_id"`
	Content          string     `json:"content"`
	UserId           int        `json:"user_id"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedBy        *int       `json:"updated_by"`
	IsLikedByUser    bool       `json:"liked"`
	IsDislikedByUser bool       `json:"disliked"`
	NumberOfLikes    int        `json:"number_of_likes"`
	NumberOfDislikes int        `json:"number_of_dislikes"`
	// //Post             Post                      `json:"post"`
	User          User `json:"user"`
	RepliesCount  int  `json:"repliesCount"`
	ISCreatedByMe bool `json:"isCreatedByMe"`
}
