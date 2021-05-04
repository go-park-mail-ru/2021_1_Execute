package delivery

import (
	"2021_1_Execute/internal/files"
	"2021_1_Execute/internal/session"
	"2021_1_Execute/internal/tasks"

	"github.com/labstack/echo"
)

type TasksHandler struct {
	sessionHD session.SessionHandler
	taskUC    tasks.TaskUsecase
	fileUT    files.FileUtil
}

func NewTasksHandler(e *echo.Echo, sessionHD session.SessionHandler, taskUC tasks.TaskUsecase, fileUT files.FileUtil) {
	handler := &TasksHandler{
		sessionHD: sessionHD,
		taskUC:    taskUC,
		fileUT:    fileUT,
	}

	e.GET("api/tasks/:id/", handler.GetTask)
	e.POST("api/tasks/", handler.PostTask)
	e.PATCH("api/tasks/:id/", handler.PatchTask)
	e.DELETE("api/tasks/:id/", handler.DeleteTask)
	e.DELETE("/api/attachments/:id/", handler.DeleteAttachment)
}
