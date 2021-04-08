package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

type sessionRepository struct {
}

func NewSessionRepository() domain.SessionsRepository {
	return &sessionRepository{}
}

func (sessionUC *sessionRepository) IsAuthorized(ctx context.Context, uuid domain.UUID) (domain.User, error) {
	//todo
	return domain.User{}, nil
}
func (sessionUC *sessionRepository) SetSession(ctx context.Context, userID int) (domain.UUID, error) {
	//todo
	return "", nil
}
func (sessionUC *sessionRepository) DeleteSession(ctx context.Context, uuid domain.UUID) error {
	//todo
	return nil
}
