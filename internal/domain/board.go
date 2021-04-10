package domain

import "context"

type Board struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Row struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BoardRepository interface {
	AddBoard(Board)
}

type BoardUsecase interface {
	AddBoard(ctx context.Context, board Board) (int, error)
	AddRow(ctx context.Context, row Row, boardID int) (int, error)
	AddTask(ctx context.Context, task Task, rowID int) (int, error)
	AddOwner(ctx context.Context, boardID int, userID int) error
	UpdateBoard(ctx context.Context, board Board) error
	UpdateRow(ctx context.Context, row Row) error
	UpdateTask(ctx context.Context, task Task) error
	DeleteBoard(ctx context.Context, boardID int) error
	DeleteRow(ctx context.Context, rowID int) error
	DeleteTask(ctx context.Context, taskID int) error
}
