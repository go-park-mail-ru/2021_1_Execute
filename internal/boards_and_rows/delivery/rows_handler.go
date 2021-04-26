package delivery

import (
	"2021_1_Execute/internal/boards_and_rows"
	http_models "2021_1_Execute/internal/boards_and_rows/models/http"
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (handler *BoardsHandler) PostRow(c echo.Context) error {
	input := new(http_models.PostRowRequest)
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

	rowID, err := handler.boardUC.AddRow(context.Background(), http_models.PostRowToRow(input), input.BoardID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, http_models.PostRowResponce{ID: rowID})
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

	return c.JSON(http.StatusOK, http_models.GetRowResponce{Row: http_models.FullRowInfoToBoardRow(row)})
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

func (handler *BoardsHandler) PatchRow(c echo.Context) error {
	input := new(http_models.PatchRowRequest)
	input.CarryOver = http_models.MoveObject{-1, -1}
	input.Move = http_models.MoveObject{-1, -1}
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
		newRow := boards_and_rows.Row{
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

		task, err := handler.taskUC.GetTask(context.Background(), input.Move.CardID, userID)
		if err != nil {
			return err
		}

		err = handler.boardUC.UpdateTasksPositions(context.Background(), rowID, task.ID, input.Move.NewPosition, userID)
		if err != nil {
			return err
		}

	}
	if input.CarryOver.NewPosition >= 0 || input.CarryOver.CardID >= 0 {
		if !(input.CarryOver.NewPosition >= 0 && input.CarryOver.CardID >= 0) {
			return errors.Wrap(domain.BadRequestError, "Need new position and ID")
		}

		err = handler.taskUC.CarryOver(context.Background(), input.CarryOver.CardID, rowID, input.CarryOver.NewPosition, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
