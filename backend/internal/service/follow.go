package service

import (
	"backend/internal/repository"
	"log"
	"net/http"
)

func FollowRequest(followerID, followedID int) (int, int) {
	if followerID == followedID {
		return 0, http.StatusBadRequest
	}

	isPublic, statusCode := repository.ProfilePrivacyByUserId(followedID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		return 0, statusCode
	}
	var frID int
	var err error

	if !isPublic {
		frID, err = repository.SendFollowRequest(followerID, followedID)
		if err != nil {
			log.Println("Error inserting/updating follow request:", err)
			return 0, http.StatusInternalServerError
		}
	} else {
		return 0, http.StatusBadRequest
	}

	return frID, http.StatusOK
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
