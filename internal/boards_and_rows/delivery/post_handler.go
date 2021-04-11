package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type PostBoardRequest struct {
	Name string `json:"name"`
}
type PostBoardResponce struct {
	ID int `json:"id"`
}

func (handler *BoardsHandler) PostBoard(c echo.Context) error {
	input := new(PostBoardRequest)
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	newBoard := domain.Board{Name: input.Name}

	ctx := context.Background()
	boardID, err := handler.boardUC.AddBoard(ctx, newBoard, userID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, PostBoardResponce{ID: boardID})
	return nil
}
