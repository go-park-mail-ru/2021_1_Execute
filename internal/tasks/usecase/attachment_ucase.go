package usecase

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"
)

func (uc *tasksUsecase) AddAttachment(ctx context.Context, taskID int, attachment tasks.Attachment, requesterID int) (int, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return -1, err
	}

	attachmentID, err := uc.tasksRepo.AddAttachment(ctx, taskID, attachment)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return attachmentID, nil
}

func (uc *tasksUsecase) DeleteAttachment(ctx context.Context, attachmentID, requesterID int) error {
	taskID, err := uc.tasksRepo.GetAttachmentTaskID(ctx, attachmentID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	err = uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.DeleteAttachment(ctx, attachmentID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) GetAttachment(ctx context.Context, attachmentID, requesterID int) (tasks.Attachment, error) {
	taskID, err := uc.tasksRepo.GetAttachmentTaskID(ctx, attachmentID)
	if err != nil {
		return tasks.Attachment{}, domain.DBErrorToServerError(err)
	}

	err = uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return tasks.Attachment{}, err
	}

	attachment, err := uc.tasksRepo.GetAttachment(ctx, attachmentID)
	if err != nil {
		return tasks.Attachment{}, domain.DBErrorToServerError(err)
	}
	if attachment.Path == "" || attachment.Name == "" {
		return tasks.Attachment{}, domain.ServerNotFoundError
	}

	return attachment, nil
}
