package usecase

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/boards_and_rows/models"
	"2021_1_Execute/internal/domain"
	"context"
)

func (uc *boardsUsecase) AddRow(ctx context.Context, row boards_and_rows.Row, boardID int, requesterID int) (int, error) {
	_, err := uc.checkAccessToBoard(ctx, boardID, requesterID)
	if err != nil {
		return 0, err
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
		return domain.DBErrorToServerError(err)
	}

	_, err = uc.checkAccessToBoard(ctx, boardID, requesterID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *boardsUsecase) GetFullRowInfo(ctx context.Context, rowID int, requesterID int) (models.FullRowInfo, error) {
	err := uc.checkRights(ctx, rowID, requesterID)
	if err != nil {
		return models.FullRowInfo{}, err
	}

	row, err := uc.boardsRepo.GetRow(ctx, rowID)
	if err != nil {
		return models.FullRowInfo{}, domain.DBErrorToServerError(err)
	}
	return uc.getFullRowInfo(ctx, row)
}

func (uc *boardsUsecase) getFullRowInfo(ctx context.Context, row boards_and_rows.Row) (models.FullRowInfo, error) {
	tasks, err := uc.boardsRepo.GetRowsTasks(ctx, row.ID)
	if err != nil {
		return models.FullRowInfo{}, domain.DBErrorToServerError(err)
	}
	return models.FullRowInfo{
		ID:       row.ID,
		Name:     row.Name,
		Position: row.Position,
		Tasks:    tasks,
	}, nil
}

func (uc *boardsUsecase) DeleteRow(ctx context.Context, rowID int, requesterID int) error {
	err := uc.checkRights(ctx, rowID, requesterID)
	if err != nil {
		return err
	}

	boardID, err := uc.boardsRepo.GetRowsBoardID(ctx, rowID)
	if err != nil {
		return err
	}

	rows, err := uc.boardsRepo.GetBoardsRows(ctx, boardID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	newPosition := len(rows) - 1

	err = uc.MoveRow(ctx, boardID, rowID, newPosition, requesterID)
	if err != nil {
		return err
	}

	err = uc.boardsRepo.DeleteRow(ctx, rowID)
	return domain.DBErrorToServerError(err)
}

func (uc *boardsUsecase) MoveRow(ctx context.Context, boardID int, rowID int, newPosition int, requesterID int) error {
	_, err := uc.checkAccessToBoard(ctx, boardID, requesterID)
	if err != nil {
		return err
	}

	rows, err := uc.boardsRepo.GetBoardsRows(ctx, boardID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	if newPosition >= len(rows) {
		newPosition = len(rows) - 1
	}

	currentRow, err := uc.boardsRepo.GetRow(ctx, rowID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	for _, row := range rows {
		if currentRow.Position > newPosition && row.Position >= newPosition && row.Position < currentRow.Position {
			row.Position += 1
			err = uc.boardsRepo.UpdateRow(ctx, row)
			if err != nil {
				return domain.DBErrorToServerError(err)
			}
		}
		if currentRow.Position < newPosition && row.Position > currentRow.Position && row.Position <= newPosition {
			row.Position -= 1
			err = uc.boardsRepo.UpdateRow(ctx, row)
			if err != nil {
				return domain.DBErrorToServerError(err)
			}
		}
	}

	currentRow.Position = newPosition
	err = uc.boardsRepo.UpdateRow(ctx, currentRow)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	return nil
}

func (uc *boardsUsecase) UpdateRow(ctx context.Context, row boards_and_rows.Row, requesterID int) error {
	err := uc.checkRights(ctx, row.ID, requesterID)
	if err != nil {
		return err
	}

	err = uc.boardsRepo.UpdateRow(ctx, row)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	return nil
}

func (uc *boardsUsecase) UpdateTasksPositions(ctx context.Context, rowID, taskID, newPos, requesterID int) error {
	err := uc.checkRights(ctx, rowID, requesterID)
	if err != nil {
		return err
	}

	err = uc.taskUC.MoveTask(ctx, taskID, newPos, requesterID)
	if err != nil {
		return err
	}

	return nil
}
