package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func SetCookie(c echo.Context, userID uint64) error {
	cookie := new(http.Cookie)
	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Error while create UUID")
	}
	sessionToken := sessionUUID.String()
	db := c.(*Database)
	(*db.Sessions)[sessionToken] = userID

	cookie.HttpOnly = true
	cookie.Name = "trello_session"
	cookie.Value = sessionToken
	cookie.Expires = time.Now().Add(12 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return nil
}
