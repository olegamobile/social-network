package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func FollowApproval(userId, targetId int) (string, int) {
	var approval string
	err := database.DB.QueryRow(`
			SELECT approval_status FROM follow_requests
			WHERE follower_id = ? AND followed_id = ?`, userId, targetId).Scan(&approval)
	if err != nil && err != sql.ErrNoRows {
		return "", http.StatusInternalServerError
	}
	return approval, http.StatusOK
}

func SendFollowRequest(followerID, followedID int) error {
	_, err := database.DB.Exec(`
        INSERT INTO follow_requests (follower_id, followed_id, approval_status)
        VALUES (?, ?, 'pending')
        ON CONFLICT(follower_id, followed_id) DO UPDATE SET approval_status='pending'
    `, followerID, followedID)

	return err
}

func StartToFollow(followerID, followedID int) error {
	_, err := database.DB.Exec(`
        INSERT INTO follow_requests (follower_id, followed_id, approval_status)
        VALUES (?, ?, 'accepted')
        ON CONFLICT(follower_id, followed_id) DO UPDATE SET approval_status='accepted'
    `, followerID, followedID)

	return err
}

func GetFollowersByUserID(userID int) ([]model.User, error) {
	query := `
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM follow_requests fr
        JOIN users u ON fr.follower_id = u.id
        WHERE fr.followed_id = ? AND fr.approval_status = 'accepted'
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &nickname); err != nil {
			return nil, err
		}
		u.Username = ""
		if nickname.Valid {
			u.Username = nickname.String
		}
		users = append(users, u)
	}
	return users, nil
}

func GetFollowedUsersByUserID(userID int) ([]model.User, error) {
	query := `
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM follow_requests fr
        JOIN users u ON fr.followed_id = u.id
        WHERE fr.follower_id = ? AND fr.approval_status = 'accepted'
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &nickname); err != nil {
			return nil, err
		}
		u.Username = ""
		if nickname.Valid {
			u.Username = nickname.String
		}
		users = append(users, u)
	}
	return users, nil
}

func GetFollowRequestsSentByUser(userID int) ([]model.User, error) {
	query := `
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM follow_requests fr
        JOIN users u ON fr.followed_id = u.id
        WHERE fr.follower_id = ? AND fr.approval_status = 'pending'
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		fmt.Println("query err at GetFollowRequestsSentByUser", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &nickname); err != nil {
			fmt.Println("scan err at GetFollowRequestsSentByUser", err)
			return nil, err
		}
		u.Username = ""
		if nickname.Valid {
			u.Username = nickname.String
		}
		users = append(users, u)
	}
	return users, nil
}

func GetFollowRequestsReceivedByUser(userID int) ([]model.User, error) {
	query := `
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM follow_requests fr
        JOIN users u ON fr.follower_id = u.id
        WHERE fr.followed_id = ? AND fr.approval_status = 'pending'
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		fmt.Println("query err at GetFollowRequestsReceivedByUser", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &nickname); err != nil {
			fmt.Println("scan err at GetFollowRequestsReceivedByUser", err)
			return nil, err
		}
		u.Username = ""
		if nickname.Valid {
			u.Username = nickname.String
		}
		users = append(users, u)
	}
	return users, nil
}

func Unfollow(followerID, followedID int) int {
	_, err := database.DB.Exec(`
        DELETE FROM follow_requests
        WHERE follower_id = ? AND followed_id = ?
    `, followerID, followedID)

	if err != nil {
		log.Println("Error deleting follow request:", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func AcceptFollowRequest(userId, followRequestId int) int {
	_, err := database.DB.Exec(`
		UPDATE follow_requests 
		SET approval_status = 'accepted'
		WHERE id = ? AND followed_id = ?
	`, followRequestId, userId)
	if err != nil {
		log.Println("Error accepting follow request:", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func DeclineFollowRequest(userId, followRequestId int) int {
	_, err := database.DB.Exec(`
		UPDATE follow_requests 
		SET approval_status = 'declined'
		WHERE id = ? AND followed_id = ?
	`, followRequestId, userId)
	if err != nil {
		log.Println("Error declining follow request:", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
