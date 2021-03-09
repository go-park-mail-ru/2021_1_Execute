package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const CookieName = "trello_session"
const CookieLifeTime = 12 * time.Hour

func SetSession(c echo.Context, userID int) error {
	cookie := new(http.Cookie)
	db := c.(*Database)

	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Error while create UUID")
	}
	sessionToken := sessionUUID.String()

	(*db.Sessions)[sessionToken] = userID

	cookie.HttpOnly = true
	cookie.Name = CookieName
	cookie.Value = sessionToken
	cookie.Expires = time.Now().Add(CookieLifeTime)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return nil
}

func DeleteSession(c echo.Context) error {
	db := c.(*Database)

	session, err := c.Cookie(CookieName)
	if err != nil {
		return &UnauthorizedError{err.Error()}
	}

	delete(*db.Sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(session)
	return nil
}
