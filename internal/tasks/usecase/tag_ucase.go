package usecase

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"
)

func (uc *tasksUsecase) AddTag(ctx context.Context, taskID int, tag tasks.Tag, requesterID int) (int, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return -1, err
	}

	tagID, err := uc.tasksRepo.AddTag(ctx, taskID, tag)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return tagID, nil
}

func (uc *tasksUsecase) AddTagToTask(ctx context.Context, taskID, tagID, requesterID int) error {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.AddTagToTask(ctx, taskID, tagID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) AddTagToBoard(ctx context.Context, boardID, tagID, requesterID int) error {
	_, err := uc.checkAccessToBoard(ctx, boardID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.AddTagToBoard(ctx, boardID, tagID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) DeleteTagFromTask(ctx context.Context, taskID, tagID, requesterID int) error {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.DeleteTagFromTask(ctx, taskID, tagID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) DeleteTagFromBoard(ctx context.Context, boardID, tagID, requesterID int) error {
	_, err := uc.checkAccessToBoard(ctx, boardID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.DeleteTagFromBoard(ctx, tagID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) DeleteTag(ctx context.Context, tagID, requesterID int) error {
	boardID, err := uc.tasksRepo.GetTagsBoardID(ctx, tagID)

	if err == nil {
		_, err := uc.checkAccessToBoard(ctx, boardID, requesterID)
		if err != nil {
			return err
		}
	} else {
		taskID, err := uc.tasksRepo.GetTagsTaskID(ctx, tagID)
		if err == nil {
			err = uc.checkRights(ctx, taskID, requesterID)
			if err != nil {
				return err
			}
		}
	}

	err = uc.tasksRepo.DeleteTag(ctx, tagID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}
