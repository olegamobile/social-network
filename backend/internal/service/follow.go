package service

import (
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
