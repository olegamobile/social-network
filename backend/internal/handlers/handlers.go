package handlers

import (
	"backend/internal/database"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	user, statusCode := service.Login(w, r)

	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := service.RemoveSession(w, r)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// "/api/me" is a common RESTful convention that returns the currently authenticated user's info
func HandleMe(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := repository.GetUserById(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		http.Error(w, "Users not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usr, statusCode := service.UserById(r)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // Error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

// SearchUsers handles the user search request.
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		json.NewEncoder(w).Encode([]model.User{}) // Return empty array for empty query
		return
	}

	users, err := repository.SearchUsers(query)
	if err != nil {
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := repository.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// handleCreatePost adds a post to the database and returns the new one
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := service.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post, statusCode := service.CreatePost(r, userID)
	if !(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices) { // error code
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("00")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		fmt.Println("01", err)
		return
	}

	// Required fields
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	firstName := strings.TrimSpace(r.FormValue("firstName"))
	lastName := strings.TrimSpace(r.FormValue("lastName"))
	dob := r.FormValue("dob")

	// Optional fields
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	about := strings.TrimSpace(r.FormValue("about"))

	if email == "" || password == "" || firstName == "" || lastName == "" || dob == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		fmt.Println("02 missing fields")
		return
	}

	// Parse and validate date of birth
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		fmt.Println("03", err)
		return
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		fmt.Println("04", err)
		return
	}

	// Handle optional avatar upload
	var avatarPath sql.NullString
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext == "" || !isAllowedImageExtension(ext) {
			http.Error(w, "Unsupported image type", http.StatusBadRequest)
			fmt.Println("05", err)
			return
		}

		filename := fmt.Sprintf("avatars/%d%s", time.Now().UnixNano(), ext)
		out, err := os.Create(filename)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			fmt.Println("06", err)
			return
		}
		defer out.Close()
		io.Copy(out, file)
		avatarPath.Valid = true
		avatarPath.String = filename
	} else {
		fmt.Println("07", err)
	}

	// Insert user into database
	query := `
	INSERT INTO users (
		email, password_hash, first_name, last_name,
		date_of_birth, avatar_path, nickname, about_me
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = database.DB.Exec(
		query,
		email, passwordHash, firstName, lastName,
		parsedDOB.Format("2006-01-02"),
		avatarPath, nullableString(nickname), nullableString(about),
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			http.Error(w, "Email already in use", http.StatusConflict)
			fmt.Println("08", err)
		} else {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			fmt.Println("09", err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

// Helper: return a sql.NullString
func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: s}
}

// Helper: check if extension is a valid image type
func isAllowedImageExtension(ext string) bool {
	ext = strings.ToLower(ext)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	return allowed[ext]
}
