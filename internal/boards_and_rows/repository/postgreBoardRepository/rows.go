package postgreBoardRepository

import (
	"2021_1_Execute/internal/domain"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreBoardRepository) AddRow(ctx context.Context, row domain.Row, boardID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "insert into rows (name, position) values ($1::text, $2::text) returning id", row.Name, row.Position)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to insert row")
	}

	var rowID int = -1

	for rows.Next() {
		err = rows.Scan(&rowID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to get row id")
		}
	}

	if rowID == -1 {
		return -1, errors.Wrap(err, "Invalid row id")
	}

	rows.Close()

	rows, err = repo.Pool.Query(ctx, "insert into boards_rows (board_id, row_id) values ($1::int, $2::int)", boardID, rowID)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to link board and row")
	}

	rows.Close()

	return rowID, nil
}

func (repo *PostgreBoardRepository) UpdateRow(ctx context.Context, row domain.Row) error {
	outdatedRow, err := repo.GetRow(ctx, row.ID)

	if err != nil {
		return errors.Wrap(err, "Unable to get outdated row")
	}

	newRow := createUpdateRowObject(outdatedRow, row)

	err = repo.updateRowQuery(ctx, newRow)

	if err != nil {
		return errors.Wrap(err, "Unable to query updating request")
	}

	return nil
}

func createUpdateRowObject(outdatedRow, newRow domain.Row) domain.Row {
	var result domain.Row

	result.ID = outdatedRow.ID

	if newRow.Name == "" {
		result.Name = outdatedRow.Name
	} else {
		result.Name = newRow.Name
	}

	if newRow.Position == -1 {
		result.Position = outdatedRow.Position
	} else {
		result.Position = newRow.Position
	}

	return result
}

func (repo *PostgreBoardRepository) updateRowQuery(ctx context.Context, row domain.Row) error {
	rows, err := repo.Pool.Query(ctx, "update rows set name = $1::text, position = $2::int where id = $3::int",
		row.Name,
		row.Position,
		row.ID,
	)

	if err != nil {
		return errors.Wrap(err, "Unable to update row")
	}

	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) GetBoardsRows(ctx context.Context, boardID int) ([]domain.Row, error) {
	rows, err := repo.Pool.Query(ctx,
		`select rows.id, rows.name, rows.position
	from rows
	inner join boards_rows as br
	on br.board_id = $1::int and br.row_id = rows.id`, boardID)

	if err != nil {
		return []domain.Row{}, errors.Wrap(err, "Unable to get boards's rows")
	}

	var boardRows []domain.Row

	for rows.Next() {
		var row domain.Row
		err = rows.Scan(&row.ID, &row.Name, &row.Position)

		if err != nil {
			return []domain.Row{}, errors.Wrap(err, "Unable to get row")
		}

		boardRows = append(boardRows, row)
	}

	rows.Close()

	return boardRows, nil
}

func (repo *PostgreBoardRepository) GetRow(ctx context.Context, rowID int) (domain.Row, error) {
	rows, err := repo.Pool.Query(ctx, "select id, name, position from rows where id = $1::int", rowID)

	if err != nil {
		return domain.Row{}, errors.Wrap(err, "Unable to get row")
	}

	var row domain.Row

	for rows.Next() {
		err = rows.Scan(&row.ID, &row.Name, &row.Position)
		if err != nil {
			return domain.Row{}, errors.Wrap(err, "Unable to read row")
		}
	}

	rows.Close()

	if row.ID == 0 {
		return domain.Row{}, domain.DBNotFoundError
	}

	return row, nil
}

func (repo *PostgreBoardRepository) DeleteRow(ctx context.Context, rowID int) error {
	row, err := repo.GetRow(ctx, rowID)
	if err != nil {
		return err
	}

	rows, err := repo.Pool.Query(ctx, "delete from rows where id = $1::int", row.ID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete row")
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

func (repo *PostgreBoardRepository) GetRowsBoardID(ctx context.Context, rowID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "select board_id from boards_rows where row_id = $1::int", rowID)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to get board id")
	}

	var boardID int = -1

	for rows.Next() {
		err = rows.Scan(&boardID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read board id")
		}
	}

	rows.Close()

	if boardID == -1 {
		return -1, domain.DBNotFoundError
	}

	return boardID, nil
}
