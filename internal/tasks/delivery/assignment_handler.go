package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func getAssignmentParams(c echo.Context) (int, int, error) {
	taskID, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		return -1, -1, domain.InternalServerError
	}
	if taskID < 0 {
		return -1, -1, domain.IDFormatError
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return -1, -1, domain.InternalServerError
	}
	if userID < 0 {
		return -1, -1, domain.IDFormatError
	}

	return taskID, userID, nil
}

func (handler *TasksHandler) PostAssignment(c echo.Context) error {
	currentUserID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, userID, err := getAssignmentParams(c)
	if err != nil {
		return err
	}

	err = handler.taskUC.Assignment(context.Background(), taskID, userID, currentUserID, "add")
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (handler *TasksHandler) DeleteAssignment(c echo.Context) error {
	currentUserID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, userID, err := getAssignmentParams(c)
	if err != nil {
		return err
	}

	err = handler.taskUC.Assignment(context.Background(), taskID, userID, currentUserID, "delete")
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
