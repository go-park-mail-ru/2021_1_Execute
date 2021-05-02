package delivery

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks/models"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (handler *TasksHandler) PostChecklist(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	input := new(models.PostChecklistRequest)

	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}
	if input.TaskID < 0 {
		return domain.IDFormatError
	}

	checklistID, err := handler.taskUC.AddChecklist(context.Background(), input.TaskID, models.PostChecklistToChecklist(input), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.PostChecklistResponse{ID: checklistID})
}

func (handler *TasksHandler) DeleteChecklist(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	checklistID, err := strconv.Atoi(c.Param("id"))
	if err != nil || checklistID < 0 {
		return domain.IDFormatError
	}

	err = handler.taskUC.DeleteChecklist(context.Background(), checklistID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (handler *TasksHandler) PatchChecklist(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	checklistID, err := strconv.Atoi(c.Param("id"))
	if err != nil || checklistID < 0 {
		return domain.IDFormatError
	}

	input := new(models.PatchChecklistRequest)

	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	err = handler.taskUC.UpdateChecklist(context.Background(), checklistID, models.PatchChecklistToChecklist(input), userID)
	if err != nil {
		return err
	}

	return nil
}
