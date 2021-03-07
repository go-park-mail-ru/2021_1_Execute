package api

import (
	"net/http"
	"time"

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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	newUser, err, code := db.CreateUser(input)
	if err != nil {
		return echo.NewHTTPError(code, err)
	}
	err = SetCookie(c, newUser.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, CreateLoginResponse(newUser))
}

func login(c echo.Context) error {
	db := c.(*Database)
	input := new(UserLoginRequest)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	flag, user := db.IsCredentialsCorect(input)
	if flag {
		err := SetCookie(c, user.ID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, CreateLoginResponse(user))
	}
	return echo.NewHTTPError(http.StatusForbidden, "Wrong pair: password, email")
}

func logout(c echo.Context) error {
	db := c.(*Database)
	session, err := c.Cookie(CookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	delete(*db.Sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(session)
	return c.NoContent(http.StatusOK)
}
