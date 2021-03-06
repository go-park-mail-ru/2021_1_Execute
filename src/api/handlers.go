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
	return c.String(http.StatusOK, "reg")
}
func login(c echo.Context) error {
	return c.String(http.StatusOK, "login")
}
func logout(c echo.Context) error {
	return c.String(http.StatusOK, "logout")
}
