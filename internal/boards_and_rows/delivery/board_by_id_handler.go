package delivery

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/boards_and_rows/models"
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (handler *BoardsHandler) GetBoardByID(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return domain.IDFormatError
	}
	if boardID < 0 {
		return domain.IDFormatError
	}

	ctx := context.Background()
	board, err := handler.boardUC.GetFullBoardInfo(ctx, boardID, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.BoardToGetResponce(board))

}

func (handler *BoardsHandler) PatchBoardByID(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil || boardID < 0 {
		return domain.IDFormatError
	}

	input := new(models.PatchBoardByIDRequest)
	input.Move.NewPosition = -1
	input.Move.RowID = -1
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err = govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	if input.Name != "" || input.Description != "" {
		err = handler.boardUC.UpdateBoard(context.Background(), boards_and_rows.Board{ID: boardID, Name: input.Name, Description: input.Description}, userID)
		if err != nil {
			return err
		}
	}
	if input.Move.RowID >= 0 || input.Move.NewPosition >= 0 {
		if !(input.Move.RowID >= 0 && input.Move.NewPosition >= 0) {
			return errors.Wrap(domain.BadRequestError, "Need rowID and new position")
		}
		err = handler.boardUC.MoveRow(context.Background(), boardID, input.Move.RowID, input.Move.NewPosition, userID)
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

func (handler *BoardsHandler) DeleteBoardByID(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil || boardID < 0 {
		return domain.IDFormatError
	}

	ctx := context.Background()
	err = handler.boardUC.DeleteBoard(ctx, boardID, userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
