package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

type boardsUsecase struct {
	boardsRepo domain.BoardRepository
}

func NewBoardsUsecase(repo domain.BoardRepository) domain.BoardUsecase {
	return &boardsUsecase{
		boardsRepo: repo,
	}
}

func (uc *boardsUsecase) AddBoard(ctx context.Context, board domain.Board, userID int) (int, error) {
	boardID, err := uc.boardsRepo.AddBoard(ctx, board)
	if err != nil {
		return 0, domain.DBErrorToServerError(err)
	}
	err = uc.boardsRepo.AddOwner(ctx, boardID, userID)
	if err != nil {
		return 0, domain.DBErrorToServerError(err)
	}
}
