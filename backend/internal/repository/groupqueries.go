package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	"net/http"
)

func ViewFullGroupOrNot(userId, targetId int) (bool, error) {
	db := database.DB
	var grId int
	err := db.QueryRow(`
	SELECT id 
	FROM group_members 
	WHERE group_id = ? AND user_id = ? AND approval_status = 'accepted' AND status = 'enable'`, targetId, userId).Scan(&grId)
	if err != nil {
		if err == sql.ErrNoRows { // No membership found so no permission to view full group
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GroupMembership(userId, targetId int) (string, int, error) {
	db := database.DB

	var approval string
	var adminId int
	err := db.QueryRow(`
	SELECT gm.approval_status, g.creator_id 
	FROM group_members gm
	JOIN groups g ON gm.group_id = g.id
	WHERE gm.group_id = ? AND gm.user_id = ? AND gm.status = 'enable' AND g.status = 'enable'`, targetId, userId).Scan(&approval, &adminId)
	if err != nil {
		if err == sql.ErrNoRows { // No membership found
			return approval, adminId, nil
		}
		return approval, adminId, err
	}
	return approval, adminId, nil
}

func GetRecommendedGroups(userID int) ([]model.Group, error) {
	query := `
		SELECT DISTINCT g.id, g.title, g.description
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.approval_status = 'accepted'
		  AND gm.user_id IN (
		      SELECT followed_id FROM follow_requests
		      WHERE follower_id = ? AND approval_status = 'accepted'
		      UNION
		      SELECT follower_id FROM follow_requests
		      WHERE followed_id = ? AND approval_status = 'accepted'
		      UNION
		      SELECT gm2.user_id
		      FROM group_members gm1
		      JOIN group_members gm2 ON gm1.group_id = gm2.group_id
		      WHERE gm1.user_id = ? AND gm1.approval_status = 'accepted'
		        AND gm2.approval_status = 'accepted'
		        AND gm2.user_id != ?
		  )
		  AND g.status = 'enable'
		  AND g.id NOT IN (
		      SELECT group_id FROM group_members
		      WHERE user_id = ? AND approval_status = 'accepted'
		  );
	`

	rows, err := database.DB.Query(query, userID, userID, userID, userID, userID)
	if err != nil {
		fmt.Println("query error at GetRecommendedGroups", err)
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description); err != nil {
			fmt.Println("scan error at GetRecommendedGroups", err)
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func SearchGroups(query string) ([]model.Group, error) {
	searchTerm := "%" + query + "%"
	sqlQuery := `
		SELECT id, title, description
		FROM groups
		WHERE status = 'enable'
		  AND (title LIKE ? OR description LIKE ?)
		ORDER BY
		  CASE
		    WHEN title LIKE ? THEN 0
		    ELSE 1
		  END,
		  title ASC;
	`

	rows, err := database.DB.Query(sqlQuery, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description); err != nil {
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func GetGroupById(groupId int) (model.Group, error) {
	var group model.Group

	err := database.DB.QueryRow(`
	SELECT id, title, description
	FROM groups
	WHERE id = ? AND status = 'enable'`, groupId).Scan(&group.ID, &group.Title, &group.Description)

	if err != nil {
		fmt.Println("error getting user by email:", err)
		return group, err
	}

	return group, nil
}

func GetGroupPostsByGroupId(groupId int) ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT gp.id, gp.user_id, gp.image_path, gp.content, gp.created_at, u.first_name, u.last_name, u.avatar_path
	FROM group_posts gp
	JOIN users u ON gp.user_id = u.id
	WHERE gp.group_id = ?
	ORDER BY gp.id DESC;`, groupId)

	if err != nil {
		fmt.Println("rows error at GetPostsByUserId", err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		var firstname, lastname string
		var avatarUrl sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &p.ImagePath, &p.Content, &p.CreatedAt, &firstname, &lastname, &avatarUrl)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}

		if avatarUrl.Valid {
			p.AvatarPath = avatarUrl.String
		} else {
			p.AvatarPath = ""
		}

		p.Username = firstname + " " + lastname
		posts = append(posts, p)
	}

	return posts, nil
}

func GetGroupMembersByGroupId(groupId int) ([]model.User, error) {
	rows, err := database.DB.Query(`
	SELECT u.id, u.first_name, u.last_name, u.avatar_path
	FROM users u
	JOIN group_members gm ON u.id = gm.user_id
	WHERE gm.group_id = ? AND gm.status = 'enable' AND u.status = 'enable' AND approval_status = 'accepted';`, groupId)

	if err != nil {
		fmt.Println("rows error at GetGroupMembersByGroupId", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var firstname, lastname string
		var avatarUrl sql.NullString

		err := rows.Scan(&u.ID, &firstname, &lastname, &avatarUrl)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}

		if avatarUrl.Valid {
			u.AvatarPath = avatarUrl.String
		} else {
			u.AvatarPath = ""
		}

		u.Username = firstname + " " + lastname
		users = append(users, u)
	}

	return users, nil
}

func GetGroupEventsByGroupId(groupId int) ([]model.Event, error) {
	rows, err := database.DB.Query(`
	SELECT e.id, g.title, e.group_id, e.title, e.description, e.event_datetime
	FROM events e
	JOIN groups g ON e.group_id = g.id
	WHERE e.group_id = ?;`, groupId)

	if err != nil {
		fmt.Println("rows error at GetGroupMembersByGroupId", err)
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		err := rows.Scan(&e.ID, &e.Group, &e.GroupID, &e.Title, &e.Description, &e.EventDate)
		if err != nil {
			fmt.Println("scan error at GetPostsByUserId", err)
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func InsertGroupPost(userID, groupID int, content string, imagePath *string) (int64, string, error) {
	query := `
		INSERT INTO group_posts (user_id, group_id, content, image_path)
		VALUES (?, ?, ?, ?)
	`
	// Use Exec for INSERT statements when you don't expect rows to be returned directly.
	result, err := database.DB.Exec(query, userID, groupID, content, imagePath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to execute insert: %w", err)
	}

	// Get the last inserted row ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get last insert ID: %w", err)
	}

	var createdAt string
	row := database.DB.QueryRow("SELECT created_at FROM group_posts WHERE id = ?", id)
	err = row.Scan(&createdAt)
	if err != nil {
		return id, "", fmt.Errorf("failed to get created_at for ID %d: %w", id, err)
	}

	return id, createdAt, nil
}

func GroupRequest(userID, groupID int) (int, int) {
	query := `
	INSERT INTO group_members (group_id, user_id, approval_status, status)
	VALUES (?, ?, 'pending', 'enable')
	ON CONFLICT(group_id, user_id) DO UPDATE SET
		approval_status = 'pending',
		status = 'enable',
		updated_by = ?
	`
	_, err := database.DB.Exec(query, groupID, userID, userID)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	var id int
	row := database.DB.QueryRow("SELECT id FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
	err = row.Scan(&id)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	return id, http.StatusOK
}

func LeaveGroup(userID, groupID int) int {
	query := `
	UPDATE group_members
	SET status = 'delete', updated_by = ?
	WHERE group_id = ? AND user_id = ?
	`
	res, err := database.DB.Exec(query, userID, groupID, userID)
	if err != nil {
		return http.StatusInternalServerError
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return http.StatusNotFound
	}
	return http.StatusOK
}

func DeleteGroup(userID, groupID int) int {
	query := `
	UPDATE groups
	SET status = 'delete', updated_by = ?
	WHERE id = ? AND creator_id = ?
	`
	res, err := database.DB.Exec(query, userID, groupID, userID)
	if err != nil {
		return http.StatusInternalServerError
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return http.StatusForbidden // Not allowed or group doesn't exist
	}
	return http.StatusOK
}

func GetAdminIdByGroupId(groupId int) (int, error) {
	var adminId int
	err := database.DB.QueryRow(`
	SELECT creator_id
	FROM groups
	WHERE id =?`, groupId).Scan(&adminId)
	if err != nil {
		return 0, err
	}
	return adminId, nil
}

func CreateGroup(group model.Group, userId int) (int, error) {
	query := `
	INSERT INTO groups (creator_id, title, description)
	VALUES (?, ?, ?)`

	res, err := database.DB.Exec(query, userId, group.Title, group.Description)
	if err != nil {
		fmt.Println("exec error creating group:", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func AddGroupMember(userId, groupId int) error {
	query := `
	INSERT INTO group_members (group_id, user_id, approval_status)
	VALUES (?, ?, 'accepted')	
	ON CONFLICT (group_id, user_id) DO UPDATE SET
		status = 'enable',
		updated_by = ?
	`
	_, err := database.DB.Exec(query, groupId, userId, userId)

	if err != nil {
		fmt.Println("exec error adding group member:", err)
	}
	return err
}

func ApproveGroupRequest(userID, groupID, adminID int, action string) int {
	query := `
		UPDATE group_members SET
			approval_status = ?,
			updated_by = ?
		WHERE 
			status = 'enable' AND 
			approval_status = 'pending' AND 
			user_id = ? AND 
			group_id = ?
	`
	_, err := database.DB.Exec(query, action, adminID, userID, groupID)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
