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

func (handler *TasksHandler) PostComment(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	input := new(models.PostCommentRequest)

	if err = c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	commentID, err := handler.taskUC.AddComment(context.Background(), models.CommentRequestToComment(input), input.TaskID, userID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.PostCommentResponse{ID: commentID})
}

func (handler *TasksHandler) DeleteComment(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil || commentID < 0 {
		return domain.IDFormatError
	}

	err = handler.taskUC.DeleteComment(context.Background(), commentID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
