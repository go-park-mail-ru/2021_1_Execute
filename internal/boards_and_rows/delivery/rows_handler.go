package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type postRowRequest struct {
	BoardID  int    `json:"board_id"`
	Name     string `json:"name" valid:"name"`
	Position int    `json:"position"`
}

type postRowResponce struct {
	ID int `json:"id"`
}

func postRowToRow(request *postRowRequest) domain.Row {
	return domain.Row{
		Name:     request.Name,
		Position: request.Position,
	}
}

func (handler *BoardsHandler) PostRow(c echo.Context) error {
	input := new(postRowRequest)
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}
	if input.BoardID < 0 {
		return domain.IDFormatError
	}
	if input.Position < 0 {
		return errors.Wrap(domain.BadRequestError, "Position should be non negative")
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	rowID, err := handler.boardUC.AddRow(context.Background(), postRowToRow(input), input.BoardID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, postRowResponce{ID: rowID})
}

type getRowResponce struct {
	Row boardRow `json:"row"`
}

type boardRow struct {
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

func fullRowInfoToBoardRow(row domain.FullRowInfo) boardRow {
	tasks := make(map[int]boardTask)
	for _, task := range row.Tasks {
		tasks[task.Position] = boardTask{
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

func (handler *BoardsHandler) GetRow(c echo.Context) error {
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil || rowID < 0 {
		return domain.IDFormatError
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	row, err := handler.boardUC.GetFullRowInfo(context.Background(), rowID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, getRowResponce{Row: fullRowInfoToBoardRow(row)})
}

func (handler *BoardsHandler) DeleteRow(c echo.Context) error {
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil || rowID < 0 {
		return domain.IDFormatError
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	err = handler.boardUC.DeleteRow(context.Background(), rowID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
