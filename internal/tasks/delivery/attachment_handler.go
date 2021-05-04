package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (handler *TasksHandler) DeleteAttachment(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	attachmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil || attachmentID < 0 {
		return domain.IDFormatError
	}

	attachment, err := handler.taskUC.GetAttachment(context.Background(), attachmentID, userID)
	if err != nil {
		return err
	}

	err = handler.fileUT.DeleteFile(attachment.Path)
	if err != nil {
		return err
	}

	err = handler.taskUC.DeleteAttachment(context.Background(), attachmentID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
