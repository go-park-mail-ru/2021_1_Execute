package http

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const CookieName = "trello_session"
const CookieLifeTime = 12 * time.Hour

type sessionHandler struct {
	sessionRepo domain.SessionsRepository
}

type SessionHandler interface {
	SetSession(c echo.Context, userID int) error
}

func NewSessionHandler(e *echo.Context, repo domain.SessionsRepository) SessionHandler {
	return &sessionHandler{
		sessionRepo: repo,
	}
}

func (handler *sessionHandler) SetSession(c echo.Context, userID int) error {

	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(domain.InternalServerError, "Error while create UUID")
	}
	sessionToken := sessionUUID.String()
	ctx := context.Background()
	err = handler.sessionRepo.SetSession(ctx, userID, sessionToken)
	if err != nil {
		return err
	}
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

func (handler *sessionHandler) DeleteSession(c echo.Context) error {

	session, err := c.Cookie(CookieName)
	if err != nil {
		return UnauthorizedError
	}

	delete(*db.Sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(session)
	return nil
}

func (handler *sessionHandler) IsAuthorized(c echo.Context) (domain.User, error) {

	_, ok := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
	}

	return c.NoContent(http.StatusOK)
}
