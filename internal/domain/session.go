package domain

import (
	"context"
)

type UUID string
type Sessions map[UUID]int

type SessionsRepository interface {
	IsAuthorized(ctx context.Context, uuid UUID) (User, error)
	SetSession(ctx context.Context, userID int) (error, UUID)
	DeleteSession(ctx context.Context, uuid UUID) error
}

type SessionsUsecase interface {
	IsAuthorized(ctx context.Context, uuid UUID) (User, error)
	SetSession(ctx context.Context, userID int) (UUID, error)
	DeleteSession(ctx context.Context, uuid UUID) error
}
