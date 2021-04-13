package delivery

import "2021_1_Execute/internal/domain"

type TasksHandler struct {
	sessionHD domain.SessionHandler
	taskUC    domain.TaskUsecase
}
