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

type FullBoardInfo struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Owner       User          `json:"-"`
	Rows        []FullRowInfo `json:"rows"`
}
type FullRowInfo struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Position int              `json:"position"`
	Tasks    []OutterTaskInfo `json:"tasks"`
}

type OutterTaskInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type BoardUsecase interface {
	AddBoard(ctx context.Context, board Board) (int, error)
	AddRow(ctx context.Context, row Row, boardID int) (int, error)
	AddTask(ctx context.Context, task Task, rowID int) (int, error)

	UpdateBoard(ctx context.Context, board Board) error
	UpdateRow(ctx context.Context, row Row) error
	UpdateTask(ctx context.Context, task Task) error

	MoveRow(ctx context.Context, boardID int, rowID int, newPosition int) error
	MoveTask(ctx context.Context, cardID int, newPosition int) error
	CarryOverTask(ctx context.Context, cardID int, newPosition int) error

	DeleteBoard(ctx context.Context, boardID int) error
	DeleteRow(ctx context.Context, rowID int) error
	DeleteTask(ctx context.Context, taskID int) error

	GetFullBoardInfo(ctx context.Context, boardID int, requesterID int) (FullBoardInfo, error)
	GetUsersBoards(ctx context.Context, userID int) ([]Board, error)
	GetTask(ctx context.Context, taskID int) (Task, error)
	GetRow(ctx context.Context, rowID int) (Row, error)
}

type BoardRepository interface {
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

	GetUsersBoards(ctx context.Context, userID int) ([]Board, error)
	GetBoardsRows(ctx context.Context, boardID int) ([]Row, error)
	GetRowsTasks(ctx context.Context, rowID int) ([]Task, error)
	GetTask(ctx context.Context, taskID int) (Task, error)
	GetRow(ctx context.Context, rowID int) (Row, error)
	GetBoard(ctx context.Context, boardID int) (Board, error)
}
