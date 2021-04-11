package delivery

import (
	"2021_1_Execute/internal/domain"

	"github.com/labstack/echo"
)

type BoardsHandler struct {
	boardUC   domain.BoardUsecase
	sessionHD domain.SessionHandler
}

func NewBoardsHandler(e *echo.Echo, boardUC domain.BoardUsecase, sessionHD domain.SessionHandler) {
	handler := &BoardsHandler{
		boardUC:   boardUC,
		sessionHD: sessionHD,
	}

	e.GET("api/boards", handler.GetUsersBoards)
	e.POST("api/boards", handler.PostBoard)
	e.GET("api/boards/:id", handler.GetBoardByID)
	e.PATCH("api/boards/:id", handler.PatchBoardByID)
	e.DELETE("api/boards/:id", handler.DeleteBoardByID)
}

func (handler *BoardsHandler) PostBoard(c echo.Context) error {
	return nil //todo
}

func (handler *BoardsHandler) PatchBoardByID(c echo.Context) error {
	return nil //todo
}
func (handler *BoardsHandler) DeleteBoardByID(c echo.Context) error {
	return nil //todo
}
