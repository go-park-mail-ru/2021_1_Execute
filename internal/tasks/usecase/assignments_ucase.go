package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

func (uc *tasksUsecase) Assignment(ctx context.Context, taskID, userID, requesterID int, typeOfAction string) error {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	switch typeOfAction {
	case "add":
		err = uc.tasksRepo.AddUserToTask(ctx, taskID, userID)
	case "delete":
		err = uc.tasksRepo.DeleteUserFromTask(ctx, taskID, userID)
	}
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}
