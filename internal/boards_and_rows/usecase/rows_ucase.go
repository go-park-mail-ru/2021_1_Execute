package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

func (uc *boardsUsecase) AddRow(ctx context.Context, row domain.Row, boardID int, requesterID int) (int, error) {
	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		return 0, domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return 0, domain.ForbiddenError
	}

	rowID, err := uc.boardsRepo.AddRow(ctx, row, boardID)
	if err != nil {
		return 0, domain.DBErrorToServerError(err)
	}
	return rowID, nil
}

func (uc *boardsUsecase) checkRights(ctx context.Context, rowID int, requesterID int) error {
	boardID, err := uc.boardsRepo.GetRowsBoardID(ctx, rowID)
	if err != nil {
		domain.DBErrorToServerError(err)
	}

	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return domain.ForbiddenError
	}
	return nil
}

func (uc *boardsUsecase) GetFullRowInfo(ctx context.Context, rowID int, requesterID int) (domain.FullRowInfo, error) {
	err := uc.checkRights(ctx, rowID, requesterID)

	row, err := uc.boardsRepo.GetRow(ctx, rowID)
	if err != nil {
		return domain.FullRowInfo{}, domain.DBErrorToServerError(err)
	}
	return uc.getFullRowInfo(ctx, row)
}

func (uc *boardsUsecase) getFullRowInfo(ctx context.Context, row domain.Row) (domain.FullRowInfo, error) {
	tasks, err := uc.boardsRepo.GetRowsTasks(ctx, row.ID)
	if err != nil {
		return domain.FullRowInfo{}, domain.DBErrorToServerError(err)
	}
	return domain.FullRowInfo{
		ID:       row.ID,
		Name:     row.Name,
		Position: row.Position,
		Tasks:    tasks,
	}, nil
}

func (uc *boardsUsecase) DeleteRow(ctx context.Context, rowID int, requesterID int) error {
	err := uc.checkRights(ctx, rowID, requesterID)

	err = uc.boardsRepo.DeleteRow(ctx, rowID)
	return domain.DBErrorToServerError(err)
}
