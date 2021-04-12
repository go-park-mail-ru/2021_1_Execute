package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

type boardsUsecase struct {
	boardsRepo domain.BoardRepository
	userUC     domain.UserUsecase
}

func NewBoardsUsecase(repo domain.BoardRepository, userUsercase domain.UserUsecase) domain.BoardUsecase {
	return &boardsUsecase{
		boardsRepo: repo,
		userUC:     userUsercase,
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
	return boardID, nil
}

func (uc *boardsUsecase) GetUsersBoards(ctx context.Context, userID int) ([]domain.Board, error) {
	boards, err := uc.boardsRepo.GetUsersBoards(ctx, userID)
	if err != nil {
		return []domain.Board{}, domain.DBErrorToServerError(err)
	}
	return boards, nil
}

func (uc *boardsUsecase) GetFullBoardInfo(ctx context.Context, boardID int, requesterID int) (domain.FullBoardInfo, error) {
	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		return domain.FullBoardInfo{}, domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return domain.FullBoardInfo{}, domain.ForbiddenError
	}

	owner, err := uc.userUC.GetUserByID(ctx, ownerID)
	if requesterID != ownerID {
		return domain.FullBoardInfo{}, domain.ForbiddenError
	}

	board, err := uc.boardsRepo.GetBoard(ctx, boardID)
	if err != nil {
		return domain.FullBoardInfo{}, domain.DBErrorToServerError(err)
	}

	rows, err := uc.boardsRepo.GetBoardsRows(ctx, boardID)
	if err != nil {
		return domain.FullBoardInfo{}, domain.DBErrorToServerError(err)
	}

	fullRowsInfo := []domain.FullRowInfo{}
	for _, row := range rows {
		tasks, err := uc.boardsRepo.GetRowsTasks(ctx, row.ID)
		if err != nil {
			return domain.FullBoardInfo{}, domain.DBErrorToServerError(err)
		}
		fullRowsInfo = append(fullRowsInfo, domain.FullRowInfo{
			ID:       row.ID,
			Name:     row.Name,
			Position: row.Position,
			Tasks:    tasks,
		})
	}
	return domain.FullBoardInfo{
		ID:          boardID,
		Name:        board.Name,
		Description: board.Description,
		Owner:       owner,
		Rows:        fullRowsInfo,
	}, nil
}

func (uc *boardsUsecase) DeleteBoard(ctx context.Context, boardID int, requesterID int) error {
	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return domain.ForbiddenError
	}

	err = uc.boardsRepo.DeleteBoard(ctx, boardID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	return nil
}
