package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"net/http"
)

func CreateEvent(event model.Event, userID int) (model.Event, error) {
	var err error

	event.CreatorID = userID

	event.Creator, err = repository.GetUserById(userID, false)
	if err != nil {
		fmt.Println("error getting creator at HandleCreateEvent:", err)
		return event, err
	}

	group, err := repository.GetGroupById(event.GroupID)
	if err != nil {
		fmt.Println("error getting group at HandleCreateEvent:", err)
		return event, err
	}
	event.Group = group.Title

	id, err := repository.CreateEvent(event)
	if err != nil {
		fmt.Println("error creating event at HandleCreateEvent:", err)
		return event, err
	}
	event.ID = &id

	return event, nil
}

func RespondToEvent(resp model.EventResponse, userID int) int {
	var err error

	resp.UserID = userID

	if resp.Response != "going" && resp.Response != "not_going" && resp.Response != "pending" {
		return http.StatusBadRequest
	}

	oldResponse, err := repository.GetEventResponse(resp.EventID, resp.UserID)
	if err != nil {
		return http.StatusBadRequest
	}

	if oldResponse == resp.Response { // remove old response when clicking same button
		resp.Response = "pending"
	}

	err = repository.SaveEventResponse(resp.EventID, resp.UserID, resp.Response)
	if err != nil {
		fmt.Println("Saving event response failed:", err)
		return http.StatusBadRequest
	}

	return http.StatusOK
}
