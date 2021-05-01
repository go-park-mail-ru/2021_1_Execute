package delivery

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (handler *BoardsHandler) getAdminRequestParams(c echo.Context) (int, int, int, error) {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return 0, 0, 0, err
	}

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, 0, 0, domain.IDFormatError
	}
	if boardID < 0 {
		return 0, 0, 0, domain.IDFormatError
	}

	newUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return 0, 0, 0, domain.IDFormatError
	}
	if newUserID < 0 {
		return 0, 0, 0, domain.IDFormatError
	}

	return boardID, userID, newUserID, nil
}

func (handler *BoardsHandler) AddAdminToBoard(c echo.Context) error {
	boardID, userID, newUserID, err := handler.getAdminRequestParams(c)

	ctx := context.Background()
	err = handler.boardUC.AddAdminToBoard(ctx, boardID, newUserID, userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (handler *BoardsHandler) DeleteAdminFromBoard(c echo.Context) error {
	boardID, userID, newUserID, err := handler.getAdminRequestParams(c)

	ctx := context.Background()
	err = handler.boardUC.DeleteAdminFromBoard(ctx, boardID, newUserID, userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
