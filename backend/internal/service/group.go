package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"net/http"
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

// get invite info: group id and user id
func ApproveGroupInvitation(inviteID, userID int, action string) int {
	// get invite info: group id and user id
	groupID, invitedID, err := repository.GetGroupInviteInfo(inviteID)
	if err != nil {
		fmt.Println("Failed to get group invitation info", err)
		return http.StatusInternalServerError
	}

	if userID != invitedID {
		fmt.Println("Only own invitations could be approved")
		return http.StatusUnauthorized
	}

	membership, err := Membership(userID, groupID)
	if err != nil {
		fmt.Println("Failed to determine group membership status", err)
		return http.StatusInternalServerError
	}

	var statusCode int
	switch membership {
	case "accepted":
		// User is already a member, just update invite status
		statusCode = repository.UpdateGroupInviteStatus(inviteID, userID, action)
	case "pending", "declined", "":
		// Update both invite and member status
		statusCode = repository.UpdateGroupInviteStatus(inviteID, userID, action)
		if action == "accepted" {
			err = repository.AddGroupMember(userID, groupID)
			if err != nil {
				fmt.Println("Failed to add group member", err)
				return http.StatusInternalServerError
			}
		}

	default:
		fmt.Println("Invalid membership status", err)
		return http.StatusInternalServerError
	}

	return statusCode
}

func GetUsersMembership(invitables []model.InvitableUser, users []model.User, groupId int) ([]model.InvitableUser, error) {
	for _, user := range users {
		var inv model.InvitableUser
		inv.User = user
		membership, err := repository.GetMembershipStatus(user.ID, groupId)
		if err != nil {
			fmt.Println("error getting membership at HandleGroupInvitationSearch:", err)
			return nil, err
		}
		inv.Membership = membership
		invitables = append(invitables, inv)
	}
	return invitables, nil
}
