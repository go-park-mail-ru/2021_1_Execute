package domain

import "context"

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Position    int    `json:"position"`
	Description string `json:"description"`
}

type TaskUsecase interface {
	AddTask(ctx context.Context, task Task, rowID int) (int, error)

	UpdateTask(ctx context.Context, task Task) error

	DeleteTask(ctx context.Context, taskID int) error

	GetTask(ctx context.Context, taskID int) (Task, error)
}

type TaskRepository interface {
	AddTask(ctx context.Context, task Task, rowID int) (int, error)

	UpdateTask(ctx context.Context, task Task) error

	DeleteTask(ctx context.Context, taskID int) error

	GetTask(ctx context.Context, taskID int) (Task, error)
}
