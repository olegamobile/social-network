package service

import "backend/internal/repository"

func Membership(userId, groupId int) (string, error) {
	approval, adminId, err := repository.GroupMembership(userId, groupId)
	if err != nil {
		return "", err
	}

	if userId == adminId {
		return "admin", nil
	}

	return approval, nil
}
