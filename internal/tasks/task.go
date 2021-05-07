package tasks

import (
	"context"
)

type Assignment struct {
	UserID int `json:"id"`
}

type Task struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Position    int          `json:"position"`
	Description string       `json:"description"`
	Comments    []Comment    `json:"comments,omitempty"`
	Assignments []Assignment `json:"users,omitempty"`
	Checklists  []Checklist  `json:"checklists,omitempty"`
}

type Comment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Author int    `json:"author"`
	Time   string `json:"time"`
}

type Field struct {
	Name string `json:"name"`
	Done bool   `json:"isDone"`
}

type Checklist struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
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

	AddComment(ctx context.Context, comment Comment, taskID, requesterID int) (int, error)
	GetComment(ctx context.Context, commentID, requesterID int) (Comment, error)
	DeleteComment(ctx context.Context, commentID, requesterID int) error
	Assignment(ctx context.Context, taskID, userID, requesterID int, typeOfAction string) error
	AddChecklist(ctx context.Context, taskID int, checklist Checklist, requesterID int) (int, error)
	DeleteChecklist(ctx context.Context, checklistID, requesterID int) error
	UpdateChecklist(ctx context.Context, checklistID int, checklist Checklist, requesterID int) error
}

type TaskRepository interface {
	AddTask(ctx context.Context, task Task, rowID int) (int, error)
	ChangeRow(ctx context.Context, taskID int, newRowID int) error

	UpdateTask(ctx context.Context, task Task) error

	DeleteTask(ctx context.Context, taskID int) error

	GetTask(ctx context.Context, taskID int) (Task, error)
	GetTasksBoardID(ctx context.Context, taskID int) (int, error)
	GetTasksRowID(ctx context.Context, taskID int) (int, error)

	AddComment(ctx context.Context, comment Comment, taskID int) (int, error)
	GetComment(ctx context.Context, commentID int) (Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
	GetCommentsTaskID(ctx context.Context, commentID int) (int, error)
	AddUserToTask(ctx context.Context, taskID, userID int) error
	DeleteUserFromTask(ctx context.Context, taskID, userID int) error
	GetTasksAssignments(ctx context.Context, taskID int) ([]int, error)
	AddChecklist(ctx context.Context, taskID int, checklist Checklist) (int, error)
	DeleteChecklist(ctx context.Context, checklistID int) error
	UpdateChecklist(ctx context.Context, checklistID int, checklist Checklist) error
	GetTasksChecklists(ctx context.Context, taskID int) ([]Checklist, error)
	GetChecklistsTaskID(ctx context.Context, checklistID int) (int, error)
}
