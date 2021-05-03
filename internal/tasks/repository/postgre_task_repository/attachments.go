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

	if attachmentID == -1 {
		return -1, errors.New("Invalid attachment's id")
	}

	return attachmentID, nil
}
