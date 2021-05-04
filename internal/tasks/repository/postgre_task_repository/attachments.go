package postgre_task_repository

import (
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreTaskRepository) AddAttachment(ctx context.Context, taskID int, attachment tasks.Attachment) (int, error) {
	rows, err := repo.Pool.Query(ctx, "insert into attachments (file_name, path, task_id) values ($1::text, $2::text, $3::int) returning id",
		attachment.Name, attachment.Path, taskID)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to insert attachment")
	}
	defer rows.Close()

	var attachmentID int = -1

	for rows.Next() {
		err = rows.Scan(&attachmentID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read attachment's id")
		}
	}

	if rows.Err() != nil {
		return -1, errors.Wrap(err, "Violated some constraints")
	}

	if attachmentID == -1 {
		return -1, errors.New("Invalid attachment's id")
	}

	return attachmentID, nil
}

func (repo *PostgreTaskRepository) DeleteAttachment(ctx context.Context, attachmentID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from attachments where id = $1::int", attachmentID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete attachment")
	}
	rows.Close()
	return nil
}

func (repo *PostgreTaskRepository) GetTasksAttachments(ctx context.Context, taskID int) ([]tasks.Attachment, error) {
	rows, err := repo.Pool.Query(ctx, "select (id, file_name, path) from attachments where task_id = $1::int", taskID)
	if err != nil {
		return []tasks.Attachment{}, errors.Wrap(err, "Unable to get task's attachments")
	}
	defer rows.Close()

	var attachments []tasks.Attachment

	for rows.Next() {
		var attachment tasks.Attachment
		err = rows.Scan(&attachment.ID, &attachment.Name, &attachment.Path)
		if err != nil {
			return []tasks.Attachment{}, errors.Wrap(err, "Unable to read task's attachment")
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func (repo *PostgreTaskRepository) GetAttachmentTaskID(ctx context.Context, attachmentID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "select task_id from attachments where id = $1::int", attachmentID)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to get attachment's task id")
	}
	defer rows.Close()

	var taskID int = -1

	for rows.Next() {
		err = rows.Scan(&taskID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read attachment's task id")
		}
	}

	if taskID == -1 {
		return -1, errors.New("Invalid task id")
	}

	return taskID, nil
}

func (repo *PostgreTaskRepository) GetAttachment(ctx context.Context, attachmentID int) (tasks.Attachment, error) {
	rows, err := repo.Pool.Query(ctx, "select (id, file_name, path) from attachments where id = $1::int", attachmentID)
	if err != nil {
		return tasks.Attachment{}, errors.Wrap(err, "Unable to get attachment")
	}
	defer rows.Close()

	var attachment tasks.Attachment

	for rows.Next() {
		err = rows.Scan(&attachment.ID, &attachment.Name, &attachment.Path)
		if err != nil {
			return tasks.Attachment{}, errors.Wrap(err, "Unable to read attachment")
		}
	}

	return attachment, nil
}
