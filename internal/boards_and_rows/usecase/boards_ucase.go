package usecase

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/boards_and_rows/models"
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"2021_1_Execute/internal/users"
	"context"
)

type boardsUsecase struct {
	boardsRepo boards_and_rows.BoardRepository
	userUC     users.UserUsecase
	taskUC     tasks.TaskUsecase
}

func NewBoardsUsecase(repo boards_and_rows.BoardRepository, userUsercase users.UserUsecase, taskUsercase tasks.TaskUsecase) boards_and_rows.BoardUsecase {
	return &boardsUsecase{
		boardsRepo: repo,
		userUC:     userUsercase,
		taskUC:     taskUsercase,
	}
}

func (uc *boardsUsecase) AddBoard(ctx context.Context, board boards_and_rows.Board, userID int) (int, error) {
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

func (uc *boardsUsecase) GetUsersBoards(ctx context.Context, userID int) ([]boards_and_rows.Board, error) {
	boards, err := uc.boardsRepo.GetUsersBoards(ctx, userID)
	if err != nil {
		return []boards_and_rows.Board{}, domain.DBErrorToServerError(err)
	}
	return boards, nil
}

func (uc *boardsUsecase) GetFullBoardInfo(ctx context.Context, boardID int, requesterID int) (models.FullBoardInfo, error) {
	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		return models.FullBoardInfo{}, domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return models.FullBoardInfo{}, domain.ForbiddenError
	}

	owner, err := uc.userUC.GetUserByID(ctx, ownerID)
	if requesterID != ownerID {
		return models.FullBoardInfo{}, domain.ForbiddenError
	}

	board, err := uc.boardsRepo.GetBoard(ctx, boardID)
	if err != nil {
		return models.FullBoardInfo{}, domain.DBErrorToServerError(err)
	}

	rows, err := uc.boardsRepo.GetBoardsRows(ctx, boardID)
	if err != nil {
		return models.FullBoardInfo{}, domain.DBErrorToServerError(err)
	}

	fullRowsInfo := []models.FullRowInfo{}
	for _, row := range rows {
		rowInfo, err := uc.getFullRowInfo(ctx, row)
		if err != nil {
			return models.FullBoardInfo{}, err
		}

		fullRowsInfo = append(fullRowsInfo, rowInfo)
	}
	return models.FullBoardInfo{
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

func (uc *boardsUsecase) UpdateBoard(ctx context.Context, board boards_and_rows.Board, requesterID int) error {
	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, board.ID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return domain.ForbiddenError
	}
	err = uc.boardsRepo.UpdateBoard(ctx, board)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	return nil
}

func (uc *boardsUsecase) changeBoardsAdmins(ctx context.Context, boardID int, newUserID int, requesterID int, isAddAction bool) error {
	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	adminsIDs, err := uc.boardsRepo.GetBoardsAdmins(ctx, boardID)

	requesterIsAdmin := false
	for id := range adminsIDs {
		if id == requesterID {
			requesterIsAdmin = true
			break
		}
	}

	if requesterID != ownerID && !requesterIsAdmin {
		return domain.ForbiddenError
	}

	if requesterID == newUserID {
		return domain.BadRequestError
	}

	_, err = uc.userUC.GetUserByID(ctx, newUserID)
	if err != nil {
		return err
	}

	if isAddAction {
		err = uc.boardsRepo.AddAdminToBoard(ctx, boardID, newUserID)
	} else {
		err = uc.boardsRepo.DeleteAdminFromBoard(ctx, boardID, newUserID)
	}
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	return nil
}

func (uc *boardsUsecase) AddAdminToBoard(ctx context.Context, boardID int, newUserID int, requesterID int) error {
	return uc.changeBoardsAdmins(ctx, boardID, newUserID, requesterID, true)
}

func (uc *boardsUsecase) DeleteAdminFromBoard(ctx context.Context, boardID int, newUserID int, requesterID int) error {
	return uc.changeBoardsAdmins(ctx, boardID, newUserID, requesterID, false)
}
