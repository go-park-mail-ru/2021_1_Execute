package delivery

import (
	"2021_1_Execute/internal/domain"

	"github.com/labstack/echo"
)

type TasksHandler struct {
	sessionHD domain.SessionHandler
	taskUC    domain.TaskUsecase
}

func NewTasksHandler(e *echo.Echo, sessionHD domain.SessionHandler, taskUC domain.TaskUsecase) {
	handler := &TasksHandler{
		sessionHD: sessionHD,
		taskUC:    taskUC,
	}

	e.GET("api/tasks/:id", handler.GetTask)
	e.POST("api/tasks", handler.PostTask)
	e.PATCH("api/tasks/:id", handler.PatchTask)
	e.DELETE("api/tasks/:id", handler.DeleteTask)
}
