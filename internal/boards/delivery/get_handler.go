package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
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

type GetBoardByIDResponce struct {
	Board getBoardByIDResponceContent `json:"board"`
}
type getBoardByIDResponceContent struct {
	ID          int               `json:"id"`
	Access      string            `json:"access"`
	IsStared    bool              `json:"isStared"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Users       boardUsers        `json:"users"`
	Rows        map[int]boardRows `json:"rows"`
}
type boardUser struct {
	ID     int    `json:"id"`
	Avatar string `json:"avatar"`
}
type boardUsers struct {
	Owner   boardUser   `json:"owner"`
	Admins  []boardUser `json:"admins"`
	Members []boardUser `json:"members"`
}
type boardRows struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Position int               `json:"position"`
	Tasks    map[int]boardTask `json:"tasks"`
}
type boardTask struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

func BoardToGetResponce(board domain.FullBoardInfo) GetBoardByIDResponce {
	boardUsers := boardUsers{
		Owner:   boardUser{ID: board.Owner.ID, Avatar: board.Owner.Avatar},
		Admins:  []boardUser{},
		Members: []boardUser{},
	}

	rows := make(map[int]boardRows)
	for _, row := range board.Rows {
		tasks := make(map[int]boardTask)
		for _, task := range row.Tasks {
			tasks[task.Position] = boardTask{
				ID:       task.ID,
				Name:     task.Name,
				Position: task.Position,
			}
		}
		rows[row.Position] = boardRows{
			ID:       row.ID,
			Name:     row.Name,
			Position: row.Position,
			Tasks:    tasks,
		}
	}

	content := getBoardByIDResponceContent{
		ID:          board.ID,
		Access:      "",
		IsStared:    false,
		Name:        board.Name,
		Description: board.Description,
		Rows:        rows,
		Users:       boardUsers,
	}
	return GetBoardByIDResponce{Board: content}
}

func BoardsToGetResponce(boards []domain.Board) GetBoardsResponce {
	responce := []getBoardsResponceContent{}
	for _, board := range boards {
		responce = append(responce, getBoardsResponceContent{
			ID:          board.ID,
			Access:      "",
			IsStared:    false,
			Name:        board.Name,
			Description: board.Description,
		})
	}
	return GetBoardsResponce{
		Boards: responce,
	}
}
func (handler *BoardsHandler) GetUsersBoards(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	boards, err := handler.boardUC.GetUsersBoards(ctx, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, BoardsToGetResponce(boards))
}

func (handler *BoardsHandler) GetBoardByID(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.Wrap(domain.ForbiddenError, "ID should be int")
	}

	ctx := context.Background()
	board, err := handler.boardUC.GetFullBoardInfo(ctx, boardID, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, BoardToGetResponce(board))

}
