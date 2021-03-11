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
	db := c.(*Database)

	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Error while create UUID")
	}
	sessionToken := sessionUUID.String()

	(*db.Sessions)[sessionToken] = userID

	cookie := http.Cookie{
		HttpOnly: true,
		Name:     CookieName,
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(CookieLifeTime),
	}
	c.SetCookie(&cookie)
	return nil
}

func DeleteSession(c echo.Context) error {
	db := c.(*Database)

	session, err := c.Cookie(CookieName)
	if err != nil {
		return UnauthorizedError
	}

	delete(*db.Sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(session)
	return nil
}
