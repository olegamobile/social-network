package model

import "database/sql"

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Birthday   string `json:"birthday"` // Birthday can be used to check if returned profile is partly hidden or not
	Password   string `json:"password"`
	About      string `json:"about_me"`
	AvatarPath string `json:"avatar_url"`
	IsPublic   bool   `json:"is_public"`
}

type Post struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	Username   string  `json:"username"`
	AvatarPath string  `json:"avatar_url"`
	Content    string  `json:"content"`
	ImagePath  *string `json:"image_path,omitempty"`
	GroupID    *int    `json:"group_id,omitempty"` // nil for regular posts
	GroupName  *string `json:"group_name,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileData struct {
	FirstName    string
	LastName     string
	DOB          string
	Nickname     string
	About        string
	AvatarPath   sql.NullString
	DeleteAvatar bool
	IsPublic     bool
}

type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type FollowRequest struct {
	TargetID int    `json:"target_id"`
	Action   string `json:"action"` // "request", "follow", "unfollow"
}

type Event struct {
	ID          *int   `json:"id,omitempty"`
	Group       string `json:"group"`
	GroupID     int    `json:"group_id"`
	CreatorID   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_datetime"`
}

type EventResponse struct {
	EventID  int    `json:"event_id"`
	UserID   int    `json:"user_id"`
	Response string `json:"response"` // going / not_going / pending
}

type GroupRequest struct {
	TargetID int    `json:"target_id"`
	Action   string `json:"action"` // "request", "leave", "delete"
}

type GroupRequestApproval struct {
	GroupID     int `json:"group_id"`
	RequesterID int `json:"requester_id"`
}

type GroupInvitation struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	GroupID string `json:"group_id"`
	Inviter int    `json:"inviter_id"`
}

type Notification struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"` // 'follow_request', 'group_invitation', 'group_join_request', 'event_creation'
	UserID        int     `json:"user_id"`
	SenderID      *int    `json:"sender_id,omitempty"`
	SenderName    *string `json:"sender_name,omitempty"`
	FollowReqID   *int    `json:"follow_req_id,omitempty"`
	GroupInviteID *int    `json:"group_invite_id,omitempty"`
	GroupID       *int    `json:"group_id,omitempty"`
	GroupTitle    *string `json:"group_title,omitempty"`
	EventID       *int    `json:"event_id,omitempty"`
	EventTitle    *string `json:"event_title,omitempty"`
	Content       *string `json:"content,omitempty"`
	IsRead        *bool   `json:"is_read,omitempty"`
	Pending       bool    `json:"pending"`
	CreatedAt     string  `json:"created_at"`
}
