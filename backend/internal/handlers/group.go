package handlers

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandleSuggestGroups(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetRecommendedGroups(userId)
	if err != nil {
		http.Error(w, "Could not fetch recommended groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func SearchGroups(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode([]model.Group{}) // Return empty array for empty query
		return
	}

	groups, err := repository.SearchGroups(query)
	if err != nil {
		http.Error(w, "Error searching groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func HandlePostsByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/posts/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	viewGroup, err := repository.ViewFullGroupOrNot(userId, targetId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var posts []model.Post

	if viewGroup {
		posts, err = repository.GetGroupPostsByGroupId(targetId)
		if err != nil {
			http.Error(w, "Failed to get posts", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func HandleMembersByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/members/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	viewGroup, err := repository.ViewFullGroupOrNot(userId, targetId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var users []model.User

	if viewGroup {
		users, err = repository.GetGroupMembersByGroupId(targetId)
		if err != nil {
			http.Error(w, "Failed to get group members", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func HandleEventsByGroupId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/events/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	viewEvents, err := repository.ViewFullGroupOrNot(userId, targetId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var events []model.Event

	if viewEvents {
		events, err = repository.GetGroupEventsByGroupId(targetId)
		if err != nil {
			http.Error(w, "Failed to get group members", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r) // replace with your actual session logic
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println("error reading data at CreateGroupPostHandler", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	groupIDStr := r.FormValue("group_id")
	if content == "" || groupIDStr == "" {
		fmt.Println("Missing fields at CreateGroupPostHandler", err)
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		fmt.Println("Invalid group_id at CreateGroupPostHandler", err)
		http.Error(w, "Invalid group_id", http.StatusBadRequest)
		return
	}

	membership, err := service.Membership(userID, groupID)
	if err != nil {
		http.Error(w, "Failed to determine group membership status", http.StatusInternalServerError)
		return
	}

	if membership != "accepted" && membership != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var imagePath *string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		savedPath, saveErr := service.SaveUploadedFile(file, header)
		if saveErr != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		imagePath = &savedPath
	} else if err != http.ErrMissingFile {
		fmt.Println("Error reading file at CreateGroupPostHandler", err)
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	id, createdAt, err := repository.InsertGroupPost(userID, groupID, content, imagePath)
	if err != nil {
		http.Error(w, "Failed to store post", http.StatusInternalServerError)
		return
	}

	user, err := repository.GetUserById(userID, true)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	group, err := repository.GetGroupById(groupID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	post := model.Post{
		ID:         int(id),
		UserID:     userID,
		Username:   user.FirstName + " " + user.LastName,
		AvatarPath: user.AvatarPath,
		Content:    content,
		ImagePath:  imagePath,
		GroupID:    &groupID,
		GroupName:  &group.Title,
		CreatedAt:  createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func HandleGroupMembership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleGroupMembership")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("ValidateSession error at HandleGroupMembership:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.GroupRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("json error at HandleGroupMembership:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var statusCode int
	switch req.Action {
	case "request":
		var gmId int
		gmId, statusCode = repository.GroupRequest(userID, req.TargetID) // 'approval_status' to pending, 'status' to enable
		if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) {
			fmt.Println("error code at HandleFollowAction:", statusCode)
			http.Error(w, http.StatusText(statusCode), statusCode)
			return
		}

		adminID, err := repository.GetAdminIdByGroupId(req.TargetID)
		if err != nil {
			fmt.Println("error getting admin in HandleGroupMembership:", err)
			http.Error(w, "error getting admin id", http.StatusBadRequest)
			return
		}

		err = repository.InsertNotification(userID, adminID, "group_join_request", gmId) // last id needs to be id at group members table
		if err != nil {
			http.Error(w, "error inserting notification in HandleGroupMembership", http.StatusBadRequest)
			return
		}
	case "leave":
		statusCode = repository.LeaveGroup(userID, req.TargetID) // 'status' to delete
	case "cancel":
		statusCode = repository.LeaveGroup(userID, req.TargetID) // 'status' to delete
		// todo: remove notification
	case "delete":
		statusCode = repository.DeleteGroup(userID, req.TargetID) // group 'status' to delete
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
	}

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleFollowAction:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func HandleGroupsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupsByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupRequests(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupRequestsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupRequestsByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch group requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupInvitations(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupInvitationsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupInvitationsByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch group invitations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupsAdministered(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupInvitationsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroupsAdministeredByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch group invitations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func HandleGroupById(w http.ResponseWriter, r *http.Request) {
	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/group/")
	targetId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	group, err := repository.GetGroupById(targetId)
	if err != nil {
		http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
		return
	}

	membership, err := service.Membership(userId, targetId)
	if err != nil {
		http.Error(w, "Failed to determine group membership status", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Group      model.Group `json:"group"`
		Membership string      `json:"membership"`
	}{
		Group:      group,
		Membership: membership,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("validate error in HandleGroupsByUserId:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var grp model.Group
	err = json.NewDecoder(r.Body).Decode(&grp)
	if err != nil {
		fmt.Println("json error at HandleCreateGroup:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	grp.ID, err = service.CreateGroup(grp, userId)
	if err != nil {
		fmt.Println("CreateGroup error in HandleCreateGroup:", err)
		http.Error(w, "failed to create group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(grp)
}

func HandleGroupRequestApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleFollowRequestApprove")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		fmt.Println("ValidateSession error at HandleFollowRequestApprove:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.GroupRequestApproval
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body at HandleFollowRequestApprove", http.StatusBadRequest)
		return
	}

	groupID := req.GroupID
	requesterID := req.RequesterID

	membership, err := service.Membership(userID, groupID)
	if err != nil {
		http.Error(w, "Failed to determine group membership status", http.StatusInternalServerError)
		return
	}

	if membership != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	action := strings.TrimPrefix(r.URL.Path, "/api/group/requests/")
	if action != "accepted" && action != "declined" {
		http.Error(w, "Invalid request action syntax", http.StatusBadRequest)
		return
	}

	statusCode := repository.ApproveGroupRequest(requesterID, groupID, userID, action)

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		fmt.Println("error code at HandleGroupRequestApprove:", statusCode)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGroupInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleFollowRequestApprove")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var groupInvite model.GroupInvitation
	err = json.NewDecoder(r.Body).Decode(&groupInvite)
	if err != nil {
		fmt.Println("json error at HandleGroupInvitation:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	groupInvite.Inviter = userID

	groupInvite.ID, err = repository.InviteToGroup(groupInvite)
	if err != nil {
		http.Error(w, "Failed to create invitation", http.StatusUnauthorized)
		return
	}

	err = repository.InsertNotification(groupInvite.Inviter, groupInvite.UserId, "group_invitation", groupInvite.ID)
	if err != nil {
		http.Error(w, "error inserting notification in HandleGroupInvitation", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGroupInvitationSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed at HandleGroupInvitationSearch")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var invitables []model.InvitableUser
	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode(invitables) // Return empty array for empty query
		return
	}

	users, err := repository.SearchUsers(query)
	if err != nil {
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	type req struct {
		GroupId string `json:"group_id"`
	}
	var groupInfo req

	err = json.NewDecoder(r.Body).Decode(&groupInfo)
	if err != nil {
		fmt.Println("json error at HandleGroupInvitationSearch:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	groupId, err := strconv.Atoi(groupInfo.GroupId)
	if err != nil {
		fmt.Println("strconv.Atoi error at HandleGroupInvitationSearch:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		var inv model.InvitableUser
		inv.User = user
		membership, err := repository.GetMembershipStatus(user.ID, groupId)
		if err != nil {
			fmt.Println("error getting membership at HandleGroupInvitationSearch:", err)
			http.Error(w, "error getting membership status", http.StatusBadRequest)
			return
		}
		inv.Membership = membership
		invitables = append(invitables, inv)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invitables)
}
