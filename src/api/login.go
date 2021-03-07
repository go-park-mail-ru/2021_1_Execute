package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func CreateLoginResponse(user User) RegistrationResponse {
	return RegistrationResponse{ID: user.ID}
}

func login(c echo.Context) error {
	db := c.(*Database)
	input := new(UserLoginRequest)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	flag, user := db.IsCredentialsCorrect(input)
	if flag {
		err := SetSession(c, user.ID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, CreateLoginResponse(user))
	}
	return echo.NewHTTPError(http.StatusForbidden, "Wrong pair: password, email")
}

func logout(c echo.Context) error {
	return DeteleSesssion(c)
}
