package postgre_task_repository

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreTaskRepository) AddTask(ctx context.Context, task tasks.Task, rowID int) (int, error) {
	repo.logger.Log(ctx, logLevelDebug, "AddTask", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.AddTask",
		"data": map[string]interface{}{
			"row_id": rowID,
			"task":   task,
		},
	})

	rows, err := repo.Pool.Query(ctx, "insert into tasks (name, description, position) values ($1::text, $2::text, $3::int) returning id", task.Name, task.Description, task.Position)

	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to insert task", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.AddTask",
			"data": map[string]interface{}{
				"row_id": rowID,
				"task":   task,
			},
			"error": err,
		})
		return -1, errors.Wrap(err, "Unable to insert task")
	}

	var taskID int = -1

	for rows.Next() {
		err = rows.Scan(&taskID)
		if err != nil {
			repo.logger.Log(ctx, logLevelError, "Unable to get task id", map[string]interface{}{
				"package": "postgre_task_repository",
				"method":  "repo.AddTask",
				"data": map[string]interface{}{
					"row_id": rowID,
					"task":   task,
				},
				"error": err,
			})
			return -1, errors.Wrap(err, "Unable to get task id")
		}
	}

	if taskID == -1 {
		return -1, errors.Wrap(err, "Invalid task id")
	}

	rows.Close()

	err = repo.connectRowAndTask(ctx, taskID, rowID)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to connect row and task")
	}

	return taskID, nil
}

func (repo *PostgreTaskRepository) UpdateTask(ctx context.Context, task tasks.Task) error {
	repo.logger.Log(ctx, logLevelDebug, "UpdateTask", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.UpdateTask",
		"data": map[string]tasks.Task{
			"task": task,
		},
	})

	outdatedTask, err := repo.GetTask(ctx, task.ID)

	if err != nil {
		return errors.Wrap(err, "Unable to get outdated task")
	}

	newTask := createUpdateTaskObject(outdatedTask, task)

	err = repo.updateTaskQuery(ctx, newTask)

	if err != nil {
		return errors.Wrap(err, "Unable to query updating request")
	}

	return nil
}

func (repo *PostgreTaskRepository) deleteConnectionBetweenTaskAndRow(ctx context.Context, taskID int) error {
	repo.logger.Log(ctx, logLevelDebug, "deleteConnectionBetweenTaskAndRow", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.deleteConnectionBetweenTaskAndRow",
		"data": map[string]int{
			"task_id": taskID,
		},
	})

	rows, err := repo.Pool.Query(ctx, "delete from rows_tasks where task_id = $1::int", taskID)
	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to delete connection between row and id", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.deleteConnectionBetweenTaskAndRow",
			"data": map[string]int{
				"task_id": taskID,
			},
			"error": err,
		})
		return errors.Wrap(err, "Unable to delete connection between row and id")
	}
	rows.Close()
	return nil
}

func createUpdateTaskObject(outdatedTask, newTask tasks.Task) tasks.Task {
	var result tasks.Task

	result.ID = outdatedTask.ID

	if newTask.Name == "" {
		result.Name = outdatedTask.Name
	} else {
		result.Name = newTask.Name
	}

	if newTask.Description == "" {
		result.Description = outdatedTask.Description
	} else {
		result.Description = newTask.Description
	}

	if newTask.Position == -1 {
		result.Position = outdatedTask.Position
	} else {
		result.Position = newTask.Position
	}

	return result
}

func (repo *PostgreTaskRepository) updateTaskQuery(ctx context.Context, task tasks.Task) error {
	repo.logger.Log(ctx, logLevelDebug, "updateTaskQuery", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.updateTaskQuery",
		"data": map[string]tasks.Task{
			"task": task,
		},
	})

	rows, err := repo.Pool.Query(ctx, "update tasks set name = $1::text, description = $2::text, position = $3::int where id = $4::int",
		task.Name,
		task.Description,
		task.Position,
		task.ID,
	)

	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to update task", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.updateTaskQuery",
			"data": map[string]tasks.Task{
				"task": task,
			},
			"error": err,
		})
		return errors.Wrap(err, "Unable to update task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) GetTask(ctx context.Context, taskID int) (tasks.Task, error) {
	repo.logger.Log(ctx, logLevelDebug, "GetTask", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.GetTask",
		"data": map[string]int{
			"task_id": taskID,
		},
	})

	rows, err := repo.Pool.Query(ctx, "select id, name, description, position from tasks where id = $1::int", taskID)

	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to get task", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.GetTask",
			"data": map[string]int{
				"task_id": taskID,
			},
			"error": err,
		})
		return tasks.Task{}, errors.Wrap(err, "Unable to get task")
	}

	var task tasks.Task

	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Position)
		if err != nil {
			repo.logger.Log(ctx, logLevelError, "Unable to read task", map[string]interface{}{
				"package": "postgre_task_repository",
				"method":  "repo.GetTask",
				"data": map[string]int{
					"task_id": taskID,
				},
				"error": err,
			})
			return tasks.Task{}, errors.Wrap(err, "Unable to read task")
		}
	}

	rows.Close()

	if task.Name == "" {
		return tasks.Task{}, domain.DBNotFoundError
	}

	return task, nil
}

func (repo *PostgreTaskRepository) DeleteTask(ctx context.Context, taskID int) error {
	repo.logger.Log(ctx, logLevelDebug, "DeleteTask", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.DeleteTask",
		"data": map[string]int{
			"task_id": taskID,
		},
	})

	task, err := repo.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Name == "" {
		return domain.DBNotFoundError
	}

	rows, err := repo.Pool.Query(ctx, "delete from tasks where id = $1::int", task.ID)
	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to delete task", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.DeleteTask",
			"data": map[string]int{
				"task_id": taskID,
			},
			"error": err,
		})
		return errors.Wrap(err, "Unable to delete task")
	}
	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) GetTasksRowID(ctx context.Context, taskID int) (int, error) {
	repo.logger.Log(ctx, logLevelDebug, "GetTasksRowID", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.GetTasksRowID",
		"data": map[string]int{
			"task_id": taskID,
		},
	})

	rows, err := repo.Pool.Query(ctx, "select row_id from rows_tasks where task_id = $1::int", taskID)
	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to get row id", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.GetTasksRowID",
			"data": map[string]int{
				"task_id": taskID,
			},
			"error": err,
		})
		return -1, errors.Wrap(err, "Unable to get row id")
	}

	var rowID int = -1

	for rows.Next() {
		err = rows.Scan(&rowID)
		if err != nil {
			repo.logger.Log(ctx, logLevelError, "Unable to read row id", map[string]interface{}{
				"package": "postgre_task_repository",
				"method":  "repo.GetTasksRowID",
				"data": map[string]int{
					"task_id": taskID,
				},
				"error": err,
			})
			return -1, errors.Wrap(err, "Unable to read row id")
		}
	}

	rows.Close()

	if rowID == -1 {
		return -1, domain.DBNotFoundError
	}

	return rowID, nil
}

func (repo *PostgreTaskRepository) GetTasksBoardID(ctx context.Context, taskID int) (int, error) {
	repo.logger.Log(ctx, logLevelDebug, "GetTasksBoardID", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.GetTasksBoardID",
		"data": map[string]int{
			"task_id": taskID,
		},
	})

	rows, err := repo.Pool.Query(ctx,
		`select br.board_id from boards_rows as br
	inner join rows_tasks as rt
	on rt.task_id = $1::int and br.row_id = rt.row_id`, taskID)
	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to get board id", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.GetTasksBoardID",
			"data": map[string]int{
				"task_id": taskID,
			},
			"error": err,
		})
		return -1, errors.Wrap(err, "Unable to get board id")
	}

	var boardID int = -1

	for rows.Next() {
		err = rows.Scan(&boardID)
		if err != nil {
			repo.logger.Log(ctx, logLevelError, "Unable to read board id", map[string]interface{}{
				"package": "postgre_task_repository",
				"method":  "repo.GetTasksBoardID",
				"data": map[string]int{
					"task_id": taskID,
				},
				"error": err,
			})
			return -1, errors.Wrap(err, "Unable to read board id")
		}
	}

	rows.Close()

	if boardID == -1 {
		return -1, domain.DBNotFoundError
	}

	return boardID, nil
}

func (repo *PostgreTaskRepository) connectRowAndTask(ctx context.Context, taskID, rowID int) error {
	repo.logger.Log(ctx, logLevelDebug, "connectRowAndTask", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.connectRowAndTask",
		"data": map[string]int{
			"task_id": taskID,
			"row_id":  rowID,
		},
	})

	rows, err := repo.Pool.Query(ctx, "insert into rows_tasks (row_id, task_id) values ($1::int, $2::int)", rowID, taskID)

	if err != nil {
		repo.logger.Log(ctx, logLevelError, "Unable to link row and task", map[string]interface{}{
			"package": "postgre_task_repository",
			"method":  "repo.connectRowAndTask",
			"data": map[string]int{
				"task_id": taskID,
				"row_id":  rowID,
			},
			"error": err,
		})
		return errors.Wrap(err, "Unable to link row and task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) ChangeRow(ctx context.Context, taskID int, newRowID int) error {
	repo.logger.Log(ctx, logLevelDebug, "ChangeRow", map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "repo.ChangeRow",
		"data": map[string]int{
			"task_id": taskID,
			"row_id":  newRowID,
		},
	})

	err := repo.deleteConnectionBetweenTaskAndRow(ctx, taskID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete outdated connections between row and task")
	}

	err = repo.connectRowAndTask(ctx, taskID, newRowID)
	if err != nil {
		return errors.Wrap(err, "Unable to change row")
	}

	return nil
}
