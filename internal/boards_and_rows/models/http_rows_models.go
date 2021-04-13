package models

import "2021_1_Execute/internal/boards_and_rows"

type PostRowRequest struct {
	BoardID  int    `json:"board_id"`
	Name     string `json:"name" valid:"name"`
	Position int    `json:"position"`
}

type PostRowResponce struct {
	ID int `json:"id"`
}

func PostRowToRow(request *PostRowRequest) boards_and_rows.Row {
	return boards_and_rows.Row{
		Name:     request.Name,
		Position: request.Position,
	}
}

type GetRowResponce struct {
	Row boardRow `json:"row"`
}

type boardRow struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Position int               `json:"position"`
	Tasks    map[int]BoardTask `json:"tasks"`
}
type BoardTask struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

func FullRowInfoToBoardRow(row boards_and_rows.FullRowInfo) boardRow {
	tasks := make(map[int]BoardTask)
	for _, task := range row.Tasks {
		tasks[task.Position] = BoardTask{
			ID:       task.ID,
			Name:     task.Name,
			Position: task.Position,
		}
	}
	return boardRow{
		ID:       row.ID,
		Name:     row.Name,
		Position: row.Position,
		Tasks:    tasks,
	}
}

type PatchRowRequest struct {
	Name      string     `json:"name,omitempty" valid:"name"`
	CarryOver MoveObject `json:"carry_over,omitempty"`
	Move      MoveObject `json:"move,omitempty"`
}

type MoveObject struct {
	CardID      int `json:"card_id"`
	NewPosition int `json:"new_position"`
}
