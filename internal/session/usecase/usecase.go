package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

type sessionsUsecase struct {
	sessionsRepo domain.SessionsRepository
}

func NewSessionUsecase(sessionsRepository domain.SessionsRepository) domain.SessionsUsecase {
	return &sessionsUsecase{
		sessionsRepo: sessionsRepository,
	}
}

func (sessionUC *sessionsUsecase) IsAuthorized(ctx context.Context, uuid domain.UUID) (domain.User, error) {
	//todo
	return domain.User{}, nil
}
func (sessionUC *sessionsUsecase) SetSession(ctx context.Context, userID int) (domain.UUID, error) {
	//todo
	return nil, ""
}
func (sessionUC *sessionsUsecase) DeleteSession(ctx context.Context, uuid domain.UUID) error {
	//todo
	return nil
}
