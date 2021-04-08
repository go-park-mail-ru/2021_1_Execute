package domain

import (
	"context"
)

type Sessions map[string]int

type SessionsRepository interface {
	IsAuthorized(ctx context.Context, uuid string) (bool, int, error)
	SetSession(ctx context.Context, userID int, uuid string) error
	DeleteSession(ctx context.Context, uuid string) error
}
