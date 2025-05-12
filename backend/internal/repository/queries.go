package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"net/http"
	"time"
)

func GetUserByEmail(req model.LoginRequest) (model.User, error) {
	var user model.User
	err := database.DB.QueryRow(`
	SELECT id, username, email, first_name, last_name, birthday, password
	FROM users WHERE email = ?`, req.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Birthday, &user.Password)

	return user, err
}

func InsertSession(sessionID string, user model.User, expiresAt time.Time) error {
	_, err := database.DB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, user.ID, expiresAt)
	return err
}

func GetUserIdAndExpirationBySessionId(cookie *http.Cookie) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM sessions WHERE id = ?`
	err := database.DB.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

func DeleteSessionById(id string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

func GetUserById(id int) (model.User, error) {
	var user model.User
	err := database.DB.QueryRow(`
		SELECT id, username, email, first_name, last_name, birthday
		FROM users WHERE id = ?`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Birthday)
	return user, err
}

func GetAllUsers() ([]model.User, error) {
	rows, _ := database.DB.Query("SELECT id, username, email, first_name, last_name, birthday FROM users")
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.Birthday)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, nil
}

func GetAllPosts() ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT posts.id, posts.user_id, users.username, posts.content, posts.created_at
	FROM posts
	JOIN users ON posts.user_id = users.id
	ORDER BY posts.id DESC;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func InsertPost(userID int, content string) (int64, string, error) {
	result, err := database.DB.Exec(
		"INSERT INTO posts (user_id, content) VALUES (?, ?)",
		userID, content,
	)

	if err != nil {
		return 0, "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	var createdAt string
	_ = database.DB.QueryRow("SELECT created_at FROM posts WHERE id = ?", id).Scan(&createdAt)

	return id, createdAt, nil
}

// SearchUsers retrieves users whose username, first_name, or last_name match the query.
func SearchUsers(query string) ([]model.User, error) {
	q := "%" + query + "%" // Add wildcards for LIKE clause
	rows, err := database.DB.Query(
		"SELECT id, username, email, first_name, last_name, birthday FROM users WHERE username LIKE ? OR first_name LIKE ? OR last_name LIKE ?",
		q, q, q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.Birthday)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
