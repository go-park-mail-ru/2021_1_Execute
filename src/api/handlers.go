package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func Router(e *echo.Echo) {
	e.POST("/api/users/", registration)
	e.POST("/api/login/", login)
	e.DELETE("/api/logout/", logout)
}

func registration(c echo.Context) error {
	db := c.(*Database)
	input := new(UserRegistrationRequest)
	if err := c.Bind(input); err != nil {
		return err
	}
	newUser, err, code := db.CreateUser(input)
	if err != nil {
		return echo.NewHTTPError(code, err.Error())
	}
	err = SetCookie(c, newUser.ID)
	if err != nil {
		return err
	}
	response := CreateRegistrationResponse(newUser)
	return c.JSON(http.StatusOK, response)
}
func login(c echo.Context) error {
	return c.String(http.StatusOK, "login")
}
func logout(c echo.Context) error {
	return c.String(http.StatusOK, "logout")
}
