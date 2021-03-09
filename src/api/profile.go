package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func createGetUserByIdResponse(user User) GetUserByIdResponse {
	return GetUserByIdResponse{
		Email:    user.Email,
		Username: user.Username,
	}
}

func GetUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	user, ok := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
	}

	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, struct {
		User GetUserByIdResponse `json:"user"`
	}{User: createGetUserByIdResponse(user)})
}

func PatchUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	input := new(PatchUserRequest)

	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	user, ok := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err = db.UpdateUser(userID, input.NewUsername, input.NewEmail, input.NewPassword)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Invalid format").Error())
	}

	return c.NoContent(http.StatusOK)
}

func DeleteUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	user, ok := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err = db.DeleteUser(userID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = DeleteSession(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
