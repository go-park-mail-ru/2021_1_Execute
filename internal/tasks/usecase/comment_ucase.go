package usecase

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"
)

func (uc *tasksUsecase) AddComment(ctx context.Context, comment tasks.Comment, taskID, requesterID int) (int, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return -1, err
	}

	commentID, err := uc.tasksRepo.AddComment(ctx, comment, taskID)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return commentID, nil
}

func (uc *tasksUsecase) GetComment(ctx context.Context, commentID, requesterID int) (tasks.Comment, error) {
	taskID, err := uc.tasksRepo.GetCommentsTaskID(ctx, commentID)
	if err != nil {
		return tasks.Comment{}, domain.DBErrorToServerError(err)
	}

	err = uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return tasks.Comment{}, err
	}

	comment, err := uc.tasksRepo.GetComment(ctx, commentID)
	if err != nil {
		return tasks.Comment{}, domain.DBErrorToServerError(err)
	}

	return comment, nil
}

func (uc *tasksUsecase) DeleteComment(ctx context.Context, commentID, requesterID int) error {
	taskID, err := uc.tasksRepo.GetCommentsTaskID(ctx, commentID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	err = uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.DeleteComment(ctx, commentID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}
