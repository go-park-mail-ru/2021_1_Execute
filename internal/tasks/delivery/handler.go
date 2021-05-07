package delivery

import (
	"2021_1_Execute/internal/session"
	"2021_1_Execute/internal/tasks"

	"github.com/labstack/echo"
)

type TasksHandler struct {
	sessionHD session.SessionHandler
	taskUC    tasks.TaskUsecase
}

func NewTasksHandler(e *echo.Echo, sessionHD session.SessionHandler, taskUC tasks.TaskUsecase) {
	handler := &TasksHandler{
		sessionHD: sessionHD,
		taskUC:    taskUC,
	}

	e.GET("api/tasks/:id/", handler.GetTask)
	e.POST("api/tasks/", handler.PostTask)
	e.PATCH("api/tasks/:id/", handler.PatchTask)
	e.DELETE("api/tasks/:id/", handler.DeleteTask)
	e.POST("api/comments/", handler.PostComment)
	e.DELETE("api/comments/:id/", handler.DeleteComment)
}
