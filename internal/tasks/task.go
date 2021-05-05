package tasks

import "context"

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Position    int    `json:"position"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags,omitempty"`
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
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

	AddTag(ctx context.Context, taskID int, tag Tag, requesterID int) (int, error)
	AddTagToTask(ctx context.Context, taskID, tagID, requesterID int) error
	AddTagToBoard(ctx context.Context, boardID, tagID, requesterID int) error
	DeleteTag(ctx context.Context, tagID, requesterID int) error
	DeleteTagFromTask(ctx context.Context, taskID, tagID, requesterID int) error
	DeleteTagFromBoard(ctx context.Context, boardID, tagID, requesterID int) error
}

type TaskRepository interface {
	AddTask(ctx context.Context, task Task, rowID int) (int, error)
	ChangeRow(ctx context.Context, taskID int, newRowID int) error

	UpdateTask(ctx context.Context, task Task) error

	DeleteTask(ctx context.Context, taskID int) error

	GetTask(ctx context.Context, taskID int) (Task, error)
	GetTasksBoardID(ctx context.Context, taskID int) (int, error)
	GetTasksRowID(ctx context.Context, taskID int) (int, error)

	AddTag(ctx context.Context, taskID int, tag Tag) (int, error)
	AddTagToTask(ctx context.Context, taskID, tagID int) error
	AddTagToBoard(ctx context.Context, boardID, tagID int) error
	DeleteTag(ctx context.Context, tagID int) error
	DeleteTagFromTask(ctx context.Context, taskID, tagID int) error
	DeleteTagFromBoard(ctx context.Context, tagID int) error
	GetTasksTags(ctx context.Context, taskID int) ([]Tag, error)
	GetBoardsTags(ctx context.Context, boardID int) ([]Tag, error)
}
