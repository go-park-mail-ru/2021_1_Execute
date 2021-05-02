package usecase

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"
)

func (uc *tasksUsecase) AddChecklist(ctx context.Context, taskID int, checklist tasks.Checklist, requesterID int) (int, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return -1, err
	}

	checklistID, err := uc.tasksRepo.AddChecklist(ctx, taskID, checklist)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return checklistID, nil
}

func (uc *tasksUsecase) DeleteChecklist(ctx context.Context, checklistID, requesterID int) error {
	taskID, err := uc.tasksRepo.GetChecklistsTaskID(ctx, checklistID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	err = uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.DeleteChecklist(ctx, checklistID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) UpdateChecklist(ctx context.Context, checklistID int, checklist tasks.Checklist, requesterID int) error {
	taskID, err := uc.tasksRepo.GetChecklistsTaskID(ctx, checklistID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	err = uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.UpdateChecklist(ctx, checklistID, checklist)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}
