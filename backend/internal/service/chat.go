package service

import (
	"backend/internal/model"
	"backend/internal/repository"
)

func SaveMessage(msg model.WSMessage) error {
	err := repository.IsFollow(msg)
	if err != nil {
		return err
	}
	err = repository.SaveMessage(msg)
	return err
}
