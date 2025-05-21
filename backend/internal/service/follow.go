package service

import (
	"backend/internal/database"
	"backend/internal/repository"
	"log"
	"net/http"
)

func FollowRequest(followerID, followedID int) int {
	if followerID == followedID {
		return http.StatusBadRequest
	}

	isPublic, statusCode := repository.ProfilePrivacyByUserId(followedID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return statusCode
	}
	if !isPublic {
		err := repository.SendFollowRequest(followerID, followedID)

		if err != nil {
			log.Println("Error inserting/updating follow request:", err)
			return http.StatusInternalServerError
		}
	} else {
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func Follow(followerID, followedID int) int {
	if followerID == followedID {
		return http.StatusBadRequest
	}

	isPublic, statusCode := repository.ProfilePrivacyByUserId(followedID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return statusCode
	}

	if isPublic {
		err := repository.StartToFollow(followerID, followedID)
		if err != nil {
			log.Println("Error inserting/updating follow:", err)
			return http.StatusInternalServerError
		}
	} else {
		return http.StatusBadRequest
	}

	return http.StatusOK
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
