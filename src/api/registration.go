package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func registration(c echo.Context) error {
	db := c.(*Database)
	input := new(UserRegistrationRequest)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser, err := db.CreateUser(input)
	if err != nil {
		return GetEchoError(err)
	}

	err = SetSession(c, newUser.ID)
	if err != nil {
		return GetEchoError(err)
	}
	return c.JSON(http.StatusOK, CreateLoginResponse(newUser))
}
