package postgre_task_repository

import (
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (repo *PostgreTaskRepository) AddChecklist(ctx context.Context, taskID int, checklist tasks.Checklist) (int, error) {
	var bitmask int64 = 0
	for i, field := range checklist.Fields {
		bitmask <<= int64(i)
		if field.Done {
			bitmask |= 1
		}
	}

	rows, err := repo.Pool.Query(ctx, "insert into checklists (task_id, name, bitmask) values ($1::int, $2::text, $3::bigint) returning id", taskID, checklist.Name, bitmask)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to create checklist")
	}

	var checklistID int = -1

	for rows.Next() {
		err = rows.Scan(&checklistID)
		if err != nil {
			rows.Close()
			return -1, errors.Wrap(err, "Unable to read checklist's id")
		}
	}
	rows.Close()

	if checklistID == -1 {
		return -1, errors.Wrap(err, "Invalid checlist's id")
	}

	for _, field := range checklist.Fields {
		rows, err := repo.Pool.Query(ctx, "update checklists set fields = array_append(fields, $1::text) where id = $2::int", field.Name, checklistID)
		if err != nil {
			return checklistID, errors.Wrap(err, "Unable to append field")
		}
		rows.Close()
	}

	return checklistID, nil
}

func (repo *PostgreTaskRepository) DeleteChecklist(ctx context.Context, checklistID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from checklists where id = $1::int", checklistID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete checklist")
	}
	rows.Close()
	return nil
}

func (repo *PostgreTaskRepository) UpdateChecklist(ctx context.Context, checklistID int, checklist tasks.Checklist) error {
	rows, err := repo.Pool.Query(ctx, "select (name) from checklists where id = $1::int", checklistID)
	if err != nil {
		return errors.Wrap(err, "Unable to get outdated checklist")
	}

	var name string

	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			rows.Close()
			return errors.Wrap(err, "Unable to read name of checklist")
		}
	}
	rows.Close()

	if name == "" {
		return errors.New("Invalid checklist")
	}

	if checklist.Name == "" {
		checklist.Name = name
	}

	var bitmask int64 = 0
	for i, field := range checklist.Fields {
		bitmask <<= int64(i)
		if field.Done {
			bitmask |= 1
		}
	}

	rows, err = repo.Pool.Query(ctx, "update checklists set name = $1::text, fields = '{}', bitmask = $2::bigint where id = $3::int",
		checklist.Name,
		bitmask,
		checklistID)
	if err != nil {
		return errors.Wrap(err, "Unable to update checklist")
	}
	rows.Close()

	for _, field := range checklist.Fields {
		rows, err := repo.Pool.Query(ctx, "update checklists set fields = array_append(fields, $1::text) where id = $2::int", field.Name, checklistID)
		if err != nil {
			return errors.Wrap(err, "Unable to update field")
		}
		rows.Close()
	}

	return nil
}

func (repo *PostgreTaskRepository) GetTasksChecklists(ctx context.Context, taskID int) ([]tasks.Checklist, error) {
	rows, err := repo.Pool.Query(ctx, "select (id, name, fields, bitmask) from checklists where task_id = $1::int", taskID)
	if err != nil {
		return []tasks.Checklist{}, errors.Wrap(err, "Unable to get task's checklists")
	}
	defer rows.Close()

	var checklists []tasks.Checklist

	for rows.Next() {
		var checklist tasks.Checklist
		var fields []string
		var bitmask int64
		err = rows.Scan(&checklist.ID, &checklist.Name, pq.Array(&fields), &bitmask)
		if err != nil {
			return []tasks.Checklist{}, errors.Wrap(err, "Unable to read checklist")
		}
		for i := len(fields) - 1; i >= 0; i-- {
			field := tasks.Field{
				Name: fields[len(fields)-i-1],
				Done: false,
			}
			if bitmask>>i&1 == 1 {
				field.Done = true
			}
			checklist.Fields = append(checklist.Fields, field)
		}
		checklists = append(checklists, checklist)
	}

	return checklists, nil
}

func (repo *PostgreTaskRepository) GetChecklistsTaskID(ctx context.Context, checklistID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "select task_id from checklists where id = $1::int", checklistID)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to get checklist's task id")
	}
	defer rows.Close()

	var taskID int = -1

	for rows.Next() {
		err = rows.Scan(&taskID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read checklist's task id")
		}
	}

	if taskID == -1 {
		return -1, errors.New("Invalid task id")
	}

	return taskID, nil
}
