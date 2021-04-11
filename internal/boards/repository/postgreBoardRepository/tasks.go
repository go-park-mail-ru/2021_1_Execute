package postgreBoardRepository

import (
	"2021_1_Execute/internal/domain"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreBoardRepository) AddTask(ctx context.Context, task domain.Task, rowID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "insert into tasks (name, description, position) values ($1::text, $2::text, $3::int) returning id", task.Name, task.Description, task.Position)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to insert task")
	}

	var taskID int = -1

	for rows.Next() {
		err = rows.Scan(&taskID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to get task id")
		}
	}

	if taskID == -1 {
		return -1, errors.Wrap(err, "Invalid task id")
	}

	rows.Close()

	rows, err = repo.Pool.Query(ctx, "insert into rows_tasks (row_id, task_id) values ($1::int, $2::int)", rowID, taskID)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to link row and task")
	}

	rows.Close()

	return taskID, nil
}

func (repo *PostgreBoardRepository) UpdateTask(ctx context.Context, task domain.Task) error {
	rows, err := repo.Pool.Query(ctx, "update tasks set name = $1::text, description = $2::text, position = $3::int where id = $3::int",
		task.Name,
		task.Description,
		task.Position,
		task.ID,
	)

	if err != nil {
		return errors.Wrap(err, "Unable to update task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) GetRowsTasks(ctx context.Context, rowID int) ([]domain.Task, error) {
	rows, err := repo.Pool.Query(ctx,
		`select tasks.id, tasks.name, tasks.description, task.position
	from tasks
	inner join rows_tasks as rt
	on rt.row_id = $1::int and rt.task_id = tasks.id`, rowID)

	if err != nil {
		return []domain.Task{}, errors.Wrap(err, "Unable to get row's tasks")
	}

	var tasks []domain.Task

	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Position)

		if err != nil {
			return []domain.Task{}, errors.Wrap(err, "Unable to get task")
		}

		tasks = append(tasks, task)
	}

	rows.Close()

	return tasks, nil
}

func (repo *PostgreBoardRepository) GetTask(ctx context.Context, taskID int) (domain.Task, error) {
	rows, err := repo.Pool.Query(ctx, "select id, name, description, position from tasks where id = $1::int", taskID)

	if err != nil {
		return domain.Task{}, errors.Wrap(err, "Unable to get task")
	}

	var task domain.Task

	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Position)
		if err != nil {
			return domain.Task{}, errors.Wrap(err, "Unable to read task")
		}
	}

	if task.Name == "" {
		return domain.Task{}, domain.DBNotFoundError
	}

	return task, nil
}

func (repo *PostgreBoardRepository) DeleteTask(ctx context.Context, taskID int) error {
	task, err := repo.GetRow(ctx, taskID)
	if err != nil {
		return err
	}

	rows, err := repo.Pool.Query(ctx, "delete from tasks where id = $1::int", task.ID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete task")
	}
	rows.Close()

	return nil
}
