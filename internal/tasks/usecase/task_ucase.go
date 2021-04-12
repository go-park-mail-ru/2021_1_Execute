package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"
)

type tasksUsecase struct {
	tasksRepo  domain.TaskRepository
	boardsRepo domain.BoardRepository
	tasksUC    domain.TaskUsecase
}

func NewTasksUsecase(taskRepo domain.TaskRepository, boardRepo domain.BoardRepository, taskUsecase domain.TaskUsecase) domain.TaskUsecase {
	return &tasksUsecase{
		tasksRepo:  taskRepo,
		boardsRepo: boardRepo,
		tasksUC:    taskUsecase,
	}
}

func (uc *tasksUsecase) AddTask(ctx context.Context, task domain.Task, rowID, requesterID int) (int, error) {
	boardID, err := uc.boardsRepo.GetRowsBoardID(ctx, rowID)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	ownerID, err := uc.boardsRepo.GetBoardsOwner(ctx, boardID)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	if requesterID != ownerID {
		return -1, domain.ForbiddenError
	}

	taskID, err := uc.tasksRepo.AddTask(ctx, task, rowID)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return taskID, nil
}

func (uc *tasksUsecase) UpdateTask(ctx context.Context, task domain.Task, requesterID int) error {
	err := uc.checkRights(ctx, task.ID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.UpdateTask(ctx, task)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) CarryOverTask(ctx context.Context, taskID, newPosition, newRowID, requesterID int) error {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.ChangeRow(ctx, taskID, newRowID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	err = uc.tasksRepo.UpdateTask(ctx, domain.Task{
		ID:       taskID,
		Position: newPosition,
	})
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) DeleteTask(ctx context.Context, taskID, requesterID int) error {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.DeleteTask(ctx, taskID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) GetTask(ctx context.Context, taskID, requesterID int) (domain.Task, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return domain.Task{}, err
	}

	task, err := uc.tasksRepo.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, domain.DBErrorToServerError(err)
	}

	return task, nil
}

func (uc *tasksUsecase) GetTasksBoardID(ctx context.Context, taskID, requesterID int) (int, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return -1, err
	}

	boardID, err := uc.tasksRepo.GetTasksBoardID(ctx, taskID)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return boardID, nil
}

func (uc *tasksUsecase) GetTasksRowID(ctx context.Context, taskID, requesterID int) (int, error) {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return -1, err
	}

	rowID, err := uc.tasksRepo.GetTasksRowID(ctx, taskID)
	if err != nil {
		return -1, domain.DBErrorToServerError(err)
	}

	return rowID, nil
}

func (uc *tasksUsecase) MoveTask(ctx context.Context, taskID, newPosition, requesterID int) error {
	err := uc.checkRights(ctx, taskID, requesterID)
	if err != nil {
		return err
	}

	err = uc.tasksRepo.UpdateTask(ctx, domain.Task{
		ID:       taskID,
		Position: newPosition,
	})
	if err != nil {
		return domain.DBErrorToServerError(err)
	}

	return nil
}

func (uc *tasksUsecase) checkRights(ctx context.Context, taskID, requesterID int) error {
	boardID, err := uc.tasksRepo.GetTasksBoardID(ctx, taskID)
	if err != nil {
		return domain.DBErrorToServerError(err)
	}
	if boardID == -1 {
		return domain.ServerNotFoundError
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
