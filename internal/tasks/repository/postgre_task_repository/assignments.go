package postgre_task_repository

import (
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreTaskRepository) AddUserToTask(ctx context.Context, taskID, userID int) error {
	rows, err := repo.Pool.Query(ctx, "insert into tasks_users (task_id, user_id) values ($1::int, $2::int)", taskID, userID)

	if err != nil {
		return errors.Wrap(err, "Unable to assign user to task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) DeleteUserFromTask(ctx context.Context, taskID, userID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from tasks_users where task_id = $1::int and user_id = $2::int", taskID, userID)

	if err != nil {
		return errors.Wrap(err, "Unable to unassign user to task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) GetTasksAssignments(ctx context.Context, taskID int) ([]int, error) {
	rows, err := repo.Pool.Query(ctx, "select user_id from tasks_users where task_id = $1::int", taskID)
	if err != nil {
		return []int{}, errors.Wrap(err, "Unable to get task's assignments")
	}
	defer rows.Close()

	var assignments []int

	for rows.Next() {
		var assignment int
		err = rows.Scan(&assignment)
		if err != nil {
			return []int{}, errors.Wrap(err, "Unable to read user ID")
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}
