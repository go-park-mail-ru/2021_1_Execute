package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const CookieName = "trello_session"

func SetSession(c echo.Context, userID uint64) error {
	cookie := new(http.Cookie)
	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Error while create UUID")
	}
	sessionToken := sessionUUID.String()
	db := c.(*Database)
	(*db.Sessions)[sessionToken] = userID

	cookie.HttpOnly = true
	cookie.Name = CookieName
	cookie.Value = sessionToken
	cookie.Expires = time.Now().Add(12 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return nil
}

func DeteleSesssion(c echo.Context) error {
	db := c.(*Database)
	session, err := c.Cookie(CookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	delete(*db.Sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(session)
	return c.NoContent(http.StatusOK)
}
