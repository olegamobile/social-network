package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
)

func SaveMessage(msg model.WSMessage) error {
	err := repository.IsFollow(msg)
	if err != nil {
		fmt.Println("error establishing follow:", err)
		return err
	}
	err = repository.SaveMessage(msg)
	return err
}

func SaveGroupMessage(msg model.WSMessage) error {
	err := repository.SaveGroupMessage(msg)
	return err
}
