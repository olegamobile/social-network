package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func UserById(r *http.Request) (model.User, int) {
	var usr model.User

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return usr, http.StatusBadRequest
	}

	usr, err = repository.GetUserById(id)
	if err != nil {
		return usr, http.StatusNotFound
	}

	return usr, http.StatusOK
}

func RegisterUser(r *http.Request) (string, int) {
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		fmt.Println("01", err)
		return "Failed to parse form", http.StatusBadRequest
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
		fmt.Println("02 missing fields")
		return "Missing required fields", http.StatusBadRequest //
	}

	// Parse and validate date of birth
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		fmt.Println("03", err)
		return "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("04", err)
		return "Failed to hash password", http.StatusInternalServerError
	}

	// Handle optional avatar upload
	var avatarPath sql.NullString
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext == "" || !isAllowedImageExtension(ext) {
			fmt.Println("05", err)
			return "Unsupported image type", http.StatusBadRequest
		}

		filename := fmt.Sprintf("avatars/%d%s", time.Now().UnixNano(), ext)
		out, err := os.Create(filename)
		if err != nil {
			fmt.Println("06", err)
			return "Failed to save image", http.StatusInternalServerError
		}
		defer out.Close()
		io.Copy(out, file)
		avatarPath.Valid = true
		avatarPath.String = filename
	} else {
		fmt.Println("07", err)
	}

	errMsg, statusCode := repository.InsertUser(passwordHash, email, firstName, lastName, parsedDOB, avatarPath, utils.NullableString(nickname), utils.NullableString(about))
	if errMsg != "" {
		return errMsg, statusCode
	}

	return "", http.StatusOK
}

// Helper: check if extension is a valid image type
func isAllowedImageExtension(ext string) bool {
	ext = strings.ToLower(ext)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	return allowed[ext]
}

func UpdateUserProfile(userID int, data model.UpdateProfileData) (model.User, error) {
	var usr model.User
	if strings.TrimSpace(data.FirstName) == "" || strings.TrimSpace(data.LastName) == "" || strings.TrimSpace(data.DOB) == "" {
		return usr, errors.New("required fields missing")
	}

	usr, err := repository.UpdateUser(userID, data)
	return usr, err
}
