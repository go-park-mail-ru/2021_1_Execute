package delivery

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

func NewSessionHandler(repo domain.SessionsRepository) domain.SessionHandler {
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
		return domain.DBErrorToServerError(err)
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
		return domain.UnauthorizedError
	}

	ctx := context.Background()
	err = handler.sessionRepo.DeleteSession(ctx, session.Value)

	if err != nil {
		return err
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(session)
	return nil
}

func (handler *sessionHandler) IsAuthorized(c echo.Context) (int, error) {
	session, err := c.Cookie(CookieName)
	if err != nil {
		return 0, domain.UnauthorizedError
	}
	ctx := context.Background()
	isAuth, userID, err := handler.sessionRepo.IsAuthorized(ctx, session.Value)

	if err != nil {
		return 0, err
	}

	if !isAuth {
		return 0, domain.UnauthorizedError
	}

	return userID, nil
}
