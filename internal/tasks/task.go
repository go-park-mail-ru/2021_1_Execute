package tasks

import "context"

type Assignment struct {
	UserID int `json:"id"`
}

type Task struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Position    int          `json:"position"`
	Description string       `json:"description"`
	Assignments []Assignment `json:"users,omitempty"`
}

type TaskUsecase interface {
	AddTask(ctx context.Context, task Task, rowID, requesterID int) (int, error)

	UpdateTask(ctx context.Context, task Task, requesterID int) error
	CarryOver(ctx context.Context, taskID, newRowID, newPosition, requesterID int) error

	DeleteTask(ctx context.Context, taskID, requesterID int) error

	GetTask(ctx context.Context, taskID, requesterID int) (Task, error)
	GetTasksBoardID(ctx context.Context, taskID, requesterID int) (int, error)
	GetTasksRowID(ctx context.Context, taskID, requesterID int) (int, error)

	MoveTask(ctx context.Context, cardID, newPosition, requesterID int) error

	Assignment(ctx context.Context, taskID, userID, requesterID int, typeOfAction string) error
}

type TaskRepository interface {
	AddTask(ctx context.Context, task Task, rowID int) (int, error)
	ChangeRow(ctx context.Context, taskID int, newRowID int) error

	UpdateTask(ctx context.Context, task Task) error

	DeleteTask(ctx context.Context, taskID int) error

	GetTask(ctx context.Context, taskID int) (Task, error)
	GetTasksBoardID(ctx context.Context, taskID int) (int, error)
	GetTasksRowID(ctx context.Context, taskID int) (int, error)

	AddUserToTask(ctx context.Context, taskID, userID int) error
	DeleteUserFromTask(ctx context.Context, taskID, userID int) error
	GetTasksAssignments(ctx context.Context, taskID int) ([]int, error)
}
