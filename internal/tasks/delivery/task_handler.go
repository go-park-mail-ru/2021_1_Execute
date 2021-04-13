package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type postTaskRequest struct {
	RowID    int    `json:"row_id"`
	Name     string `json:"name" valid:"name"`
	Position int    `json:"position"`
}

type postTaskResponse struct {
	ID int `json:"id"`
}

func taskRequestToTask(req *postTaskRequest) domain.Task {
	return domain.Task{
		Name:     req.Name,
		Position: req.Position,
	}
}

func (handler *TasksHandler) PostTask(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	input := new(postTaskRequest)

	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}
	if input.RowID < 0 {
		return domain.IDFormatError
	}
	if input.Position < 0 {
		return errors.Wrap(domain.BadRequestError, "Position should be non negative")
	}

	_, err = govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	taskID, err := handler.taskUC.AddTask(context.Background(), taskRequestToTask(input), input.RowID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, postTaskResponse{ID: taskID})
}

type getTaskResponse struct {
	Task domain.Task `json:"task"`
}

func (handler *TasksHandler) GetTask(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID < 0 {
		return domain.IDFormatError
	}

	task, err := handler.taskUC.GetTask(context.Background(), taskID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, getTaskResponse{Task: task})
}

func (handler *TasksHandler) DeleteTask(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID < 0 {
		return domain.IDFormatError
	}

	err = handler.taskUC.DeleteTask(context.Background(), taskID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

type patchTaskRequest struct {
	Name        string `json:"name,omitempty" valid:"name"`
	Description string `json:"description,omitempty" valid:"description"`
}

func patchTaskToTask(req *patchTaskRequest) domain.Task {
	return domain.Task{
		Name:        req.Name,
		Description: req.Description,
	}
}

func (handler *TasksHandler) PatchTask(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID < 0 {
		return domain.IDFormatError
	}

	input := new(patchTaskRequest)

	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err = govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	task := patchTaskToTask(input)
	task.ID = taskID

	err = handler.taskUC.UpdateTask(context.Background(), task, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
