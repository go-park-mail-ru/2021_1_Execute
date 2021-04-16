package delivery

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks/models"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (handler *TasksHandler) PostTask(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	input := new(models.PostTaskRequest)

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

	taskID, err := handler.taskUC.AddTask(context.Background(), models.TaskRequestToTask(input), input.RowID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.PostTaskResponse{ID: taskID})
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

	return c.JSON(http.StatusOK, models.GetTaskResponse{Task: task})
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

func (handler *TasksHandler) PatchTask(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID < 0 {
		return domain.IDFormatError
	}

	input := new(models.PatchTaskRequest)

	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err = govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	task := models.PatchTaskToTask(input)
	task.ID = taskID

	err = handler.taskUC.UpdateTask(context.Background(), task, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
