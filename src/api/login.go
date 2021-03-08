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

	user, isCorrect := db.IsCredentialsCorrect(input)
	if isCorrect {
		err := SetSession(c, user.ID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, CreateLoginResponse(user))
	}
	return echo.NewHTTPError(http.StatusForbidden, "Wrong pair: password, email")
}

func logout(c echo.Context) error {
	return DeteleSession(c)
}
