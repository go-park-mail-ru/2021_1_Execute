package delivery

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/session"
	"2021_1_Execute/internal/tasks"

	"github.com/labstack/echo"
)

type BoardsHandler struct {
	boardUC   boards_and_rows.BoardUsecase
	sessionHD session.SessionHandler
	taskUC    tasks.TaskUsecase
}

func NewBoardsHandler(e *echo.Echo, boardUC boards_and_rows.BoardUsecase,
	sessionHD session.SessionHandler, taskUC tasks.TaskUsecase) {

	handler := &BoardsHandler{
		boardUC:   boardUC,
		sessionHD: sessionHD,
		taskUC:    taskUC,
	}

	e.GET("api/boards/", handler.GetUsersBoards)
	e.POST("api/boards/", handler.PostBoard)

	e.GET("api/boards/:id/", handler.GetBoardByID)
	e.PATCH("api/boards/:id/", handler.PatchBoardByID)
	e.DELETE("api/boards/:id/", handler.DeleteBoardByID)

	e.POST("api/rows/", handler.PostRow)
	e.GET("api/rows/:id/", handler.GetRow)
	e.PATCH("api/rows/:id/", handler.PatchRow)
	e.DELETE("api/rows/:id/", handler.DeleteRow)
}
