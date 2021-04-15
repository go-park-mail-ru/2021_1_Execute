package boards_and_rows

import (
	"2021_1_Execute/internal/tasks"
	"2021_1_Execute/internal/users"
	"context"
)

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

type FullBoardInfo struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Owner       users.User    `json:"-"`
	Rows        []FullRowInfo `json:"rows"`
}
type FullRowInfo struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Position int          `json:"position"`
	Tasks    []tasks.Task `json:"tasks"`
}

type BoardUsecase interface {
	AddBoard(ctx context.Context, board Board, userID int) (int, error)
	AddRow(ctx context.Context, row Row, boardID int, requesterID int) (int, error)

	UpdateBoard(ctx context.Context, board Board, requesterID int) error
	UpdateRow(ctx context.Context, row Row, requesterID int) error
	MoveRow(ctx context.Context, boardID int, rowID int, newPosition int, requesterID int) error

	DeleteBoard(ctx context.Context, boardID int, requesterID int) error
	DeleteRow(ctx context.Context, rowID int, requesterID int) error

	GetFullBoardInfo(ctx context.Context, boardID int, requesterID int) (FullBoardInfo, error)
	GetUsersBoards(ctx context.Context, userID int) ([]Board, error)
	GetFullRowInfo(ctx context.Context, rowID int, requesterID int) (FullRowInfo, error)
	UpdateTasksPositions(ctx context.Context, rowID, taskID, newPos, requesterID int) error
}

type BoardRepository interface {
	AddBoard(ctx context.Context, board Board) (int, error)
	AddRow(ctx context.Context, row Row, boardID int) (int, error)
	AddOwner(ctx context.Context, boardID int, userID int) error

	UpdateBoard(ctx context.Context, board Board) error
	UpdateRow(ctx context.Context, row Row) error

	DeleteBoard(ctx context.Context, boardID int) error
	DeleteRow(ctx context.Context, rowID int) error

	GetBoard(ctx context.Context, boardID int) (Board, error)
	GetRow(ctx context.Context, rowID int) (Row, error)

	GetUsersBoards(ctx context.Context, userID int) ([]Board, error)
	GetBoardsOwner(ctx context.Context, boardID int) (int, error)

	GetBoardsRows(ctx context.Context, boardID int) ([]Row, error)
	GetRowsTasks(ctx context.Context, rowID int) ([]tasks.Task, error)

	GetRowsBoardID(ctx context.Context, rowID int) (int, error)
}
