package domain

import (
	"context"

	"github.com/labstack/echo"
)

type Sessions map[string]int

type SessionsRepository interface {
	IsAuthorized(ctx context.Context, uuid string) (bool, int, error)
	SetSession(ctx context.Context, userID int, uuid string) error
	DeleteSession(ctx context.Context, uuid string) error
}

type SessionHandler interface {
	SetSession(c echo.Context, userID int) error
	DeleteSession(c echo.Context) error
	IsAuthorized(c echo.Context) (int, error)
}
