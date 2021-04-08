package http

import (
	"net/http"

	"github.com/labstack/echo"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EntranceResponse struct {
	ID int `json:"id"`
}

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(c echo.Context) error {
	db := c.(*Database)

	input := new(UserLoginRequest)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, isCorrect := db.IsCredentialsCorrect(input)
	if isCorrect {
		err := SetSession(c, user.ID)
		if err != nil {
			return GetEchoError(err)
		}
		return c.JSON(http.StatusOK, CreateLoginResponse(user))
	}
	return echo.NewHTTPError(http.StatusForbidden, "Wrong pair: password, email")
}

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
func logout(c echo.Context) error {
	err := DeleteSession(c)
	if err != nil {
		return GetEchoError(err)
	}

	return c.NoContent(http.StatusOK)
}
