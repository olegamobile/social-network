package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
	Password  string `json:"password"`
}

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow frontend origin
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		// Allow credentials like cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// For preflight requests
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow(`
	SELECT id, username, email, first_name, last_name, birthday, password
	FROM users WHERE email = ?`, req.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Birthday, &user.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Password != req.Password { // unencrypted demo, bcrypt later
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = db.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, user.ID, expiresAt)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err == nil {
		fmt.Println("deleting cookie", cookie)
		db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	}

	fmt.Println("setting cookie to empty")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // expired
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

// "/api/me" is a common RESTful convention that returns the currently authenticated user's info
func handleMe(w http.ResponseWriter, r *http.Request) {
	userID, err := validateSession(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var user User
	err = db.QueryRow(`
		SELECT id, username, email, first_name, last_name, birthday
		FROM users WHERE id = ?`, userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Birthday)

	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	_, err := validateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, _ := db.Query("SELECT id, username, email, first_name, last_name, birthday FROM users")
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.Birthday)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	_, err := validateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var u User
	err = db.QueryRow("SELECT id, username, email, first_name, last_name, birthday FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.Birthday)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	_, err := validateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
		SELECT posts.id, posts.user_id, users.username, posts.content, posts.created_at
		FROM posts
		JOIN users ON posts.user_id = users.id
	`)
	if err != nil {
		http.Error(w, "Failed to query posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Content, &p.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to scan post row", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
