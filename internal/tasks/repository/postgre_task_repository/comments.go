package postgre_task_repository

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreTaskRepository) AddComment(ctx context.Context, comment tasks.Comment, taskID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "insert into comments (text, time, task_id, user_id) values ($1::text, timestamp $2::text, $3::int, $4::int) returning id",
		comment.Text, comment.Time.String(), taskID, comment.Author)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to create comment")
	}
	defer rows.Close()

	var commentID int = -1

	for rows.Next() {
		err = rows.Scan(&commentID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to get comment id")
		}
	}

	if rows.Err() != nil {
		return -1, errors.Wrap(rows.Err(), "")
	}
	if commentID == -1 {
		return -1, errors.Wrap(err, "Invalid comment id")
	}

	return commentID, nil
}

func (repo *PostgreTaskRepository) getTasksComments(ctx context.Context, taskID int) ([]tasks.Comment, error) {
	rows, err := repo.Pool.Query(ctx, "select (id, text, time, user_id) from comments where task_id = $1::int", taskID)

	if err != nil {
		return []tasks.Comment{}, errors.Wrap(err, "Unable to get task's comments")
	}
	defer rows.Close()

	var comments []tasks.Comment

	for rows.Next() {
		var comment tasks.Comment
		err = rows.Scan(&comment.ID, &comment.Text, &comment.Time, &comment.Author)
		if err != nil {
			return []tasks.Comment{}, errors.Wrap(err, "Unable to read")
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (repo *PostgreTaskRepository) GetComment(ctx context.Context, commentID int) (tasks.Comment, error) {
	rows, err := repo.Pool.Query(ctx, "select (id, text, time, user_id) from comments where id = $1::int", commentID)

	if err != nil {
		return tasks.Comment{}, errors.Wrap(err, "Unable to get comment")
	}
	defer rows.Close()

	var comment tasks.Comment

	for rows.Next() {
		err = rows.Scan(&comment.ID, &comment.Text, &comment.Time, &comment.Author)
		if err != nil {
			return tasks.Comment{}, errors.Wrap(err, "Unable to read comment")
		}
	}

	return comment, nil
}

func (repo *PostgreTaskRepository) DeleteComment(ctx context.Context, commentID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from comments where id = $1::int", commentID)

	if err != nil {
		return errors.Wrap(err, "Unable to delete comment")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) GetCommentsTaskID(ctx context.Context, commentID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "select task_id from comments where id = $1::int", commentID)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to get comment's taskID")
	}
	defer rows.Close()

	var taskID int = -1

	for rows.Next() {
		err = rows.Scan(&taskID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read comment's taskID")
		}
	}

	if taskID == -1 {
		return -1, domain.ServerNotFoundError
	}

	return taskID, nil
}
