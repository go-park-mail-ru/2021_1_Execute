package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
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

type patchRowRequest struct {
	Name      string     `json:"name,omitempty" valid:"name"`
	CarryOver moveObject `json:"carry_over,omitempty"`
	Move      moveObject `json:"move,omitempty"`
}

type moveObject struct {
	CardID      int `json:"card_id"`
	NewPosition int `json:"new_position"`
}

func (handler *BoardsHandler) PatchRow(c echo.Context) error {
	input := new(patchRowRequest)
	input.CarryOver = moveObject{-1, -1}
	input.Move = moveObject{-1, -1}
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil || rowID < 0 {
		return domain.IDFormatError
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	if input.Name != "" {
		newRow := domain.Row{
			ID:   rowID,
			Name: input.Name,
		}
		err = handler.boardUC.UpdateRow(context.Background(), newRow, userID)
		if err != nil {
			return err
		}
	}
	if input.Move.NewPosition >= 0 || input.Move.CardID >= 0 {
		if !(input.Move.NewPosition >= 0 && input.Move.CardID >= 0) {
			return errors.Wrap(domain.BadRequestError, "Need new position and ID")
		}
		err = handler.taskUC.MoveTask(context.Background(), input.Move.CardID, input.Move.NewPosition, userID)
		if err != nil {
			return err
		}
	}
	if input.CarryOver.NewPosition >= 0 || input.CarryOver.CardID >= 0 {
		if !(input.CarryOver.NewPosition >= 0 && input.CarryOver.CardID >= 0) {
			return errors.Wrap(domain.BadRequestError, "Need new position and ID")
		}
		err = handler.taskUC.CarryOverTask(context.Background(), input.CarryOver.CardID, input.CarryOver.NewPosition, input.CarryOver.NewPosition, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
