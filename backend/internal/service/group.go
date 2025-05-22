package service

import (
	"backend/internal/model"
	"backend/internal/repository"
)

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

func CreateGroup(group model.Group, userId int) (int, error) {
	var err error
	group.ID, err = repository.CreateGroup(group, userId)
	if err != nil {
		return 0, err
	}
	err = repository.AddGroupMember(userId, group.ID)
	return group.ID, err
}
