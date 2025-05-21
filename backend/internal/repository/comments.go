package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"fmt"
)

func ReadAllCommentsForPost(postID int, userID int) ([]model.Comment, error) {

	var comments []model.Comment
	commentMap := make(map[int]*model.Comment)
	// Updated query to include comment_id
	selectQuery := `
		SELECT 
			u.id AS user_id, u.first_name ,u.last_name ,u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by,
			c.id AS comment_id, c.post_id AS comment_post_id, c.user_id AS comment_user_id, c.content AS comment_description, 
			c.status AS comment_status, c.created_at AS comment_created_at, 
			 c.updated_by AS comment_updated_by
		FROM comments c
			INNER JOIN users u
				ON c.user_id = u.id AND c.status != 'delete' AND u.status != 'delete' AND c.post_id = ?
		ORDER BY comment_created_at asc;
	`
	rows, selectError := database.DB.Query(selectQuery, postID) // Query the database
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over rows and populate the slice
	for rows.Next() {
		var comment model.Comment
		var user model.User
		err := rows.Scan(
			// Map user fields
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,

			// Map comment fields
			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedBy,
		)
		comment.User = user
		if err != nil {
			return nil, err
		}
		if user.ID == userID {
			comment.ISCreatedByMe = true
		}

		// Handle likes/dislikes aggregation
		if _, found := commentMap[comment.ID]; !found {
			commentMap[comment.ID] = &comment
		}
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert the map of comments into a slice
	for _, comment := range commentMap {
		comments = append(comments, *comment)
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
