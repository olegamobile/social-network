package repository

import (
	"backend/internal/database"
	"backend/internal/model"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// username > nickname
// birthday > date_of_birth
// null to

func GetUserByEmail(req model.LoginRequest) (model.User, error) {
	var user model.User
	var nickname sql.NullString
	var about sql.NullString

	err := database.DB.QueryRow(`
	SELECT id, nickname, email, first_name, last_name, date_of_birth, password_hash, about_me
	FROM users WHERE email = ?`, req.Email).
		Scan(&user.ID, &nickname, &user.Email, &user.FirstName, &user.LastName, &user.Birthday, &user.Password, &about)

	if err != nil {
		fmt.Println("error getting user by email:", err)
	}

	if nickname.Valid {
		user.Username = nickname.String
	} else {
		user.Username = ""
	}

	if about.Valid {
		user.About = about.String
	} else {
		user.About = ""
	}

	return user, err
}

func InsertSession(sessionID string, user model.User, expiresAt time.Time) error {
	_, err := database.DB.Exec("INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)", sessionID, user.ID, expiresAt)
	return err
}

func GetUserIdAndExpirationBySessionId(cookie *http.Cookie) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM sessions WHERE session_token = ?`
	err := database.DB.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

func DeleteSessionById(id string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE session_token = ?", id)
	return err
}

func GetUserById(id int) (model.User, error) {
	var user model.User
	var nickname sql.NullString
	var about sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, nickname, email, first_name, last_name, date_of_birth, about_me
		FROM users WHERE id = ?`, id).
		Scan(&user.ID, &nickname, &user.Email, &user.FirstName, &user.LastName, &user.Birthday, &about)

	if nickname.Valid {
		user.Username = nickname.String
	} else {
		user.Username = ""
	}

	if about.Valid {
		user.About = about.String
	} else {
		user.About = ""
	}

	return user, err
}

func GetAllUsers() ([]model.User, error) {
	rows, _ := database.DB.Query("SELECT id, nickname, email, first_name, last_name, date_of_birth, about_me FROM users")
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		var about sql.NullString

		err := rows.Scan(&u.ID, &nickname, &u.Email, &u.FirstName, &u.LastName, &u.Birthday, &about)
		if err != nil {
			return users, err
		}
		if nickname.Valid {
			u.Username = nickname.String
		} else {
			u.Username = ""
		}

		if about.Valid {
			u.About = about.String
		} else {
			u.About = ""
		}

		users = append(users, u)
	}

	return users, nil
}

func GetAllPosts() ([]model.Post, error) {
	rows, err := database.DB.Query(`
	SELECT posts.id, posts.user_id, users.first_name, users.last_name, posts.content, posts.created_at
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
		var firstname, lastname string
		err := rows.Scan(&p.ID, &p.UserID, &firstname, &lastname, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		p.Username = firstname + " " + lastname
		posts = append(posts, p)
	}

	return posts, nil
}

func InsertPost(userID int, content string, privacy string) (int, string, error) {
	result, err := database.DB.Exec(
		"INSERT INTO posts (user_id, content, privacy_level) VALUES (?, ?, ?)",
		userID, content, privacy,
	)

	if err != nil {
		fmt.Println("error 1 at insert post", err)
		return 0, "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error 2 at insert post", err)
		return 0, "", err
	}

	var createdAt string
	_ = database.DB.QueryRow("SELECT created_at FROM posts WHERE id = ?", id).Scan(&createdAt)

	return int(id), createdAt, nil
}

// SearchUsers retrieves users whose username, first_name, or last_name match the query.
func SearchUsers(query string) ([]model.User, error) {
	q := "%" + query + "%" // Add wildcards for LIKE clause
	rows, err := database.DB.Query(`
		SELECT id, nickname, email, first_name, last_name, date_of_birth, about_me
		FROM users 
		WHERE nickname LIKE ? OR first_name LIKE ? OR last_name LIKE ?
		`,
		q, q, q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var nickname sql.NullString
		var about sql.NullString

		err := rows.Scan(&u.ID, &nickname, &u.Email, &u.FirstName, &u.LastName, &u.Birthday, &about)
		if err != nil {
			return nil, err
		}
		if nickname.Valid {
			u.Username = nickname.String
		} else {
			u.Username = ""
		}

		if about.Valid {
			u.About = about.String
		} else {
			u.About = ""
		}

		users = append(users, u)
	}

	return users, nil
}
