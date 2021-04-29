package delivery

import (
	"2021_1_Execute/internal/boards_and_rows"
	http_models "2021_1_Execute/internal/boards_and_rows/models/http"

	"2021_1_Execute/internal/domain"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (handler *BoardsHandler) GetUsersBoards(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	boards, err := handler.boardUC.GetUsersBoards(ctx, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, http_models.BoardsToGetResponce(boards))
}

func (handler *BoardsHandler) PostBoard(c echo.Context) error {
	input := new(http_models.PostBoardRequest)
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	newBoard := boards_and_rows.Board{Name: input.Name}

	ctx := context.Background()
	boardID, err := handler.boardUC.AddBoard(ctx, newBoard, userID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, http_models.PostBoardResponce{ID: boardID})
	return nil
}
