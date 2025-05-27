package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"fmt"
)

func ReadAllCommentsForPost(postID int, userID int) ([]model.Comment, error) {
	var comments []model.Comment

	selectQuery := `
SELECT 
    u.id AS user_id,
    u.first_name,
    u.last_name,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    u.updated_by AS user_updated_by,
    c.id AS comment_id,
    c.post_id AS comment_post_id,
    c.user_id AS comment_user_id,
    c.content AS comment_description, 
    c.status AS comment_status,
    c.created_at AS comment_created_at, 
    c.updated_by AS comment_updated_by
FROM comments c
INNER JOIN users u ON c.user_id = u.id
WHERE 
    c.status != 'delete' 
    AND u.status != 'delete'
    AND c.post_id = ?
ORDER BY c.created_at DESC;
`

	rows, err := database.DB.Query(selectQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,

			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		comment.User = user
		comment.ISCreatedByMe = (user.ID == userID)

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func InsertComment(content string, userID int, postID int) error {

	insertQuery := `INSERT INTO comments (post_id, content, user_id) VALUES (?, ?, ?);`
	_, insertErr := database.DB.Exec(insertQuery, postID, content, userID)

	if insertErr != nil {
		fmt.Println("Error inserting the comment", insertErr)
		return insertErr
	}

	return nil
}

func ReadAllCommentsForGroupPost(postID int, userID int) ([]model.Comment, error) {
	var comments []model.Comment

	selectQuery := `
SELECT 
    u.id AS user_id,
    u.first_name,
    u.last_name,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    u.updated_by AS user_updated_by,
    c.id AS comment_id,
    c.group_post_id AS comment_post_id,
    c.user_id AS comment_user_id,
    c.content AS comment_description, 
    c.status AS comment_status,
    c.created_at AS comment_created_at, 
    c.updated_by AS comment_updated_by
FROM group_comments c
INNER JOIN users u ON c.user_id = u.id
WHERE 
    c.status != 'delete' 
    AND u.status != 'delete'
    AND c.group_post_id = ?
ORDER BY c.created_at DESC;
`

	rows, err := database.DB.Query(selectQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,

			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		comment.User = user
		comment.ISCreatedByMe = (user.ID == userID)

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func InsertGroupComment(content string, userID int, postID int) error {

	insertQuery := `INSERT INTO group_comments (group_post_id, content, user_id) VALUES (?, ?, ?);`
	_, insertErr := database.DB.Exec(insertQuery, postID, content, userID)

	if insertErr != nil {
		fmt.Println("Error inserting the comment", insertErr)
		return insertErr
	}

	return nil
}
