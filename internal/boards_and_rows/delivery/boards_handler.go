package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type GetBoardsResponce struct {
	Boards []getBoardsResponceContent `json:"boards"`
}
type getBoardsResponceContent struct {
	ID          int    `json:"id"`
	Access      string `json:"access"`
	IsStared    bool   `json:"isStared"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func BoardsToGetResponce(boards []domain.Board) GetBoardsResponce {
	responce := []getBoardsResponceContent{}
	for _, board := range boards {
		responce = append(responce, getBoardsResponceContent{
			ID:          board.ID,
			Access:      "",
			IsStared:    false,
			Name:        board.Name,
			Description: board.Description,
		})
	}
	return GetBoardsResponce{
		Boards: responce,
	}
}
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
	return c.JSON(http.StatusOK, BoardsToGetResponce(boards))
}

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
