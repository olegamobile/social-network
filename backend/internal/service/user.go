package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
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
		avatarPath, err = uploadAvatar(file, header)
		if err != nil {
			return "Failed to save image", http.StatusInternalServerError
		}
	} else {
		fmt.Println("Image reading error at registering:", err)
	}

	errMsg, statusCode := repository.InsertUser(passwordHash, email, firstName, lastName, parsedDOB, avatarPath, utils.NullableString(nickname), utils.NullableString(about))
	if errMsg != "" {
		return errMsg, statusCode
	}

	return "", http.StatusOK
}

func UpdateUserProfile(userID int, r *http.Request) (model.User, string, int) {
	var usr model.User

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return usr, "Error parsing form", http.StatusBadRequest
	}

	updateData := model.UpdateProfileData{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		DOB:       r.FormValue("dob"),
		Nickname:  r.FormValue("nickname"),
		About:     r.FormValue("about"),
	}

	if updateData.Nickname == "null" {
		updateData.Nickname = ""
	}

	if updateData.About == "null" {
		updateData.About = ""
	}

	// Handle optional avatar
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		updateData.AvatarPath, err = uploadAvatar(file, header)
		if err != nil {
			return usr, "Failed to save image", http.StatusInternalServerError
		}
		//fmt.Println("Form avatar found at:", updateData.AvatarPath.String)
	} else {
		fmt.Println("No avatar file found at updating profile:", err)
	}

	// Handle delete_avatar flag
	if r.FormValue("delete_avatar") == "true" {
		updateData.DeleteAvatar = true
	}

	if strings.TrimSpace(updateData.FirstName) == "" || strings.TrimSpace(updateData.LastName) == "" || strings.TrimSpace(updateData.DOB) == "" {
		return usr, "required fields missing", http.StatusBadRequest
	}

	usr, errMsg, statusCode := repository.UpdateUser(userID, updateData)
	if errMsg != "" {
		return usr, errMsg, statusCode
	}

	//fmt.Println("Updated user:", usr)

	return usr, "", http.StatusOK
}

func uploadAvatar(file multipart.File, header *multipart.FileHeader) (sql.NullString, error) {
	var avatarPath sql.NullString

	ext := filepath.Ext(header.Filename)
	if ext == "" || !utils.IsAllowedImageExtension(ext) {
		fmt.Println("bad extension:", ext)
		return avatarPath, fmt.Errorf("illegal extension")
	}
	filename := fmt.Sprintf("avatars/%d%s", time.Now().UnixNano(), ext)

	dst, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating path at uploadAvatar:", err)
		return avatarPath, err
	}
	defer dst.Close()

	io.Copy(dst, file)

	avatarPath.Valid = true
	avatarPath.String = filename
	//fmt.Println("Avatar uploaded succesfully")

	return avatarPath, nil
}
