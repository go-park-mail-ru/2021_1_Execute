package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func createGetUserByIdResponse(user User) GetUserByIdResponse {
	return GetUserByIdResponse{
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.Avatar,
	}
}

func createGetUserByIdBody(user User) GetUserByIdBody {
	return GetUserByIdBody{
		Response: createGetUserByIdResponse(user),
	}
}

func GetCurrentUser(c echo.Context) error {
	db := c.(*Database)

	user, ok := db.IsAuthorized(c)
	log.Println(user, ok)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
	}

	return c.JSON(http.StatusOK, createGetUserByIdBody(user))
}

func GetUserByID(c echo.Context) error {
	fmt.Println(c.Param("id"))
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

func PatchUser(c echo.Context) error {
	input := new(PatchUserRequest)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	user, ok := db.IsAuthorized(c)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	err := db.UpdateUser(user.ID, input.NewUsername, input.NewEmail, input.NewPassword, "")

	if err != nil {
		return GetEchoError(err)
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
