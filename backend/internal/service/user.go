package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"net/http"
	"strconv"
	"strings"
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
