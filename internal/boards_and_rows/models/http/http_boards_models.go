package http

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/boards_and_rows/models"
)

type GetBoardsResponce struct {
	Boards []getBoardsResponceContent `json:"boards"`
}
type getBoardsResponceContent struct {
	ID          int    `json:"id"`
	Access      string `json:"access"`
	IsStared    bool   `json:"isStared"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func BoardsToGetResponce(boards []boards_and_rows.Board) GetBoardsResponce {
	responce := []getBoardsResponceContent{}
	for _, board := range boards {
		responce = append(responce, getBoardsResponceContent{
			ID:          board.ID,
			Access:      "guest",
			IsStared:    false,
			Name:        board.Name,
			Description: board.Description,
		})
	}
	return GetBoardsResponce{
		Boards: responce,
	}
}

type PostBoardRequest struct {
	Name string `json:"name"`
}
type PostBoardResponce struct {
	ID int `json:"id"`
}

type GetBoardByIDResponce struct {
	Board getBoardByIDResponceContent `json:"board"`
}
type getBoardByIDResponceContent struct {
	ID          int              `json:"id"`
	Access      string           `json:"access"`
	IsStared    bool             `json:"isStared"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Users       boardUsers       `json:"users"`
	Rows        map[int]boardRow `json:"rows"`
}
type boardUser struct {
	ID     int    `json:"id"`
	Avatar string `json:"avatar" validate:"url"`
}
type boardUsers struct {
	Owner   boardUser   `json:"owner,omitempty"`
	Admins  []boardUser `json:"admins,omitempty"`
	Members []boardUser `json:"members,omitempty"`
}

type PatchBoardByIDRequest struct {
	Access      string     `json:"access,omitempty"`
	IsStared    bool       `json:"isStared,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Users       boardUsers `json:"users,omitempty"`
	Move        rowsMove   `json:"move,omitempty"`
}
type rowsMove struct {
	RowID       int `json:"row_id"`
	NewPosition int `json:"new_position"`
}

func BoardToGetResponce(board models.FullBoardInfo) GetBoardByIDResponce {
	var admins []boardUser
	for _, usr := range board.Admins {
		admins = append(admins, boardUser{ID: usr.ID, Avatar: usr.Avatar})
	}

	boardUsers := boardUsers{
		Owner:   boardUser{ID: board.Owner.ID, Avatar: board.Owner.Avatar},
		Admins:  admins,
		Members: []boardUser{},
	}

	rows := make(map[int]boardRow)
	for _, row := range board.Rows {
		rows[row.Position] = FullRowInfoToBoardRow(row)
	}

	content := getBoardByIDResponceContent{
		ID:          board.ID,
		Access:      "guest",
		IsStared:    false,
		Name:        board.Name,
		Description: board.Description,
		Rows:        rows,
		Users:       boardUsers,
	}
	return GetBoardByIDResponce{Board: content}
}
