package postgre_board_repository

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (repo *PostgreBoardRepository) AddRow(ctx context.Context, row boards_and_rows.Row, boardID int) (int, error) {
	repo.log(ctx, pgx.LogLevelDebug, "AddRow", "AddRow", map[string]interface{}{
		"board_id": boardID,
		"row":      row,
	}, nil)

	rows, err := repo.Pool.Query(ctx, "insert into rows (name, position) values ($1::text, $2::int) returning id", row.Name, row.Position)

	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to insert row", "AddRow", map[string]interface{}{
			"board_id": boardID,
			"row":      row,
		}, err)
		return -1, errors.Wrap(err, "Unable to insert row")
	}

	var rowID int = -1

	for rows.Next() {
		err = rows.Scan(&rowID)
		if err != nil {
			repo.log(ctx, pgx.LogLevelError, "Unable to get row id", "AddRow", map[string]interface{}{
				"board_id": boardID,
				"row":      row,
			}, err)
			return -1, errors.Wrap(err, "Unable to get row id")
		}
	}

	if rowID == -1 {
		return -1, errors.Wrap(err, "Invalid row id")
	}

	rows.Close()

	rows, err = repo.Pool.Query(ctx, "insert into boards_rows (board_id, row_id) values ($1::int, $2::int)", boardID, rowID)

	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to link board and row", "AddRow", map[string]interface{}{
			"board_id": boardID,
			"row":      row,
		}, err)
		return -1, errors.Wrap(err, "Unable to link board and row")
	}

	rows.Close()

	return rowID, nil
}

func (repo *PostgreBoardRepository) UpdateRow(ctx context.Context, row boards_and_rows.Row) error {
	repo.log(ctx, pgx.LogLevelDebug, "UpdateRow", "UpdateRow", map[string]interface{}{
		"row": row,
	}, nil)

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

func createUpdateRowObject(outdatedRow, newRow boards_and_rows.Row) boards_and_rows.Row {
	var result boards_and_rows.Row

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

func (repo *PostgreBoardRepository) updateRowQuery(ctx context.Context, row boards_and_rows.Row) error {
	repo.log(ctx, pgx.LogLevelDebug, "updateRowQuery", "updateRowQuery", map[string]interface{}{
		"row": row,
	}, nil)

	rows, err := repo.Pool.Query(ctx, "update rows set name = $1::text, position = $2::int where id = $3::int",
		row.Name,
		row.Position,
		row.ID,
	)

	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to update row", "updateRowQuery", map[string]interface{}{
			"row": row,
		}, err)
		return errors.Wrap(err, "Unable to update row")
	}

	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) GetBoardsRows(ctx context.Context, boardID int) ([]boards_and_rows.Row, error) {
	repo.log(ctx, pgx.LogLevelDebug, "GetBoardsRows", "GetBoardsRows", map[string]interface{}{
		"board_id": boardID,
	}, nil)

	rows, err := repo.Pool.Query(ctx,
		`select rows.id, rows.name, rows.position
	from rows
	inner join boards_rows as br
	on br.board_id = $1::int and br.row_id = rows.id`, boardID)

	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to get boards's rows", "GetBoardsRows", map[string]interface{}{
			"board_id": boardID,
		}, err)
		return []boards_and_rows.Row{}, errors.Wrap(err, "Unable to get boards's rows")
	}

	var boardRows []boards_and_rows.Row

	for rows.Next() {
		var row boards_and_rows.Row
		err = rows.Scan(&row.ID, &row.Name, &row.Position)

		if err != nil {
			repo.log(ctx, pgx.LogLevelError, "Unable to get row", "GetBoardsRows", map[string]interface{}{
				"board_id": boardID,
			}, err)
			return []boards_and_rows.Row{}, errors.Wrap(err, "Unable to get row")
		}

		boardRows = append(boardRows, row)
	}

	rows.Close()

	return boardRows, nil
}

func (repo *PostgreBoardRepository) GetRow(ctx context.Context, rowID int) (boards_and_rows.Row, error) {
	repo.log(ctx, pgx.LogLevelDebug, "GetRow", "GetRow", map[string]interface{}{
		"row_id": rowID,
	}, nil)

	rows, err := repo.Pool.Query(ctx, "select id, name, position from rows where id = $1::int", rowID)

	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to get row", "GetRow", map[string]interface{}{
			"row_id": rowID,
		}, err)
		return boards_and_rows.Row{}, errors.Wrap(err, "Unable to get row")
	}

	var row boards_and_rows.Row

	for rows.Next() {
		err = rows.Scan(&row.ID, &row.Name, &row.Position)
		if err != nil {
			repo.log(ctx, pgx.LogLevelError, "Unable to read row", "GetRow", map[string]interface{}{
				"row_id": rowID,
			}, err)
			return boards_and_rows.Row{}, errors.Wrap(err, "Unable to read row")
		}
	}

	rows.Close()

	if row.ID == 0 {
		return boards_and_rows.Row{}, domain.DBNotFoundError
	}

	return row, nil
}

func (repo *PostgreBoardRepository) DeleteRow(ctx context.Context, rowID int) error {
	repo.log(ctx, pgx.LogLevelDebug, "DeleteRow", "DeleteRow", map[string]interface{}{
		"row_id": rowID,
	}, nil)

	row, err := repo.GetRow(ctx, rowID)
	if err != nil {
		return err
	}

	rows, err := repo.Pool.Query(ctx, "delete from rows where id = $1::int", row.ID)
	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to delete row", "DeleteRow", map[string]interface{}{
			"row_id": rowID,
		}, err)
		return errors.Wrap(err, "Unable to delete row")
	}
	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) GetRowsTasks(ctx context.Context, rowID int) ([]tasks.Task, error) {
	repo.log(ctx, pgx.LogLevelDebug, "GetRowsTasks", "GetRowsTasks", map[string]interface{}{
		"row_id": rowID,
	}, nil)

	rows, err := repo.Pool.Query(ctx,
		`select tasks.id, tasks.name, tasks.description, tasks.position
	from tasks
	inner join rows_tasks as rt
	on rt.row_id = $1::int and rt.task_id = tasks.id`, rowID)

	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to get row's tasks", "GetRowsTasks", map[string]interface{}{
			"row_id": rowID,
		}, err)
		return []tasks.Task{}, errors.Wrap(err, "Unable to get row's tasks")
	}

	var result []tasks.Task

	for rows.Next() {
		var task tasks.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Position)

		if err != nil {
			repo.log(ctx, pgx.LogLevelError, "Unable to read task", "GetRowsTasks", map[string]interface{}{
				"row_id": rowID,
			}, err)
			return []tasks.Task{}, errors.Wrap(err, "Unable to read task")
		}

		result = append(result, task)
	}

	rows.Close()

	return result, nil
}

func (repo *PostgreBoardRepository) GetRowsBoardID(ctx context.Context, rowID int) (int, error) {
	repo.log(ctx, pgx.LogLevelDebug, "GetRowsBoardID", "GetRowsBoardID", map[string]interface{}{
		"row_id": rowID,
	}, nil)

	rows, err := repo.Pool.Query(ctx, "select board_id from boards_rows where row_id = $1::int", rowID)
	if err != nil {
		repo.log(ctx, pgx.LogLevelError, "Unable to get board id", "GetRowsBoardID", map[string]interface{}{
			"row_id": rowID,
		}, err)
		return -1, errors.Wrap(err, "Unable to get board id")
	}

	var boardID int = -1

	for rows.Next() {
		err = rows.Scan(&boardID)
		if err != nil {
			repo.log(ctx, pgx.LogLevelError, "Unable to read board id", "GetRowsBoardID", map[string]interface{}{
				"row_id": rowID,
			}, err)
			return -1, errors.Wrap(err, "Unable to read board id")
		}
	}

	rows.Close()

	if boardID == -1 {
		return -1, domain.DBNotFoundError
	}

	return boardID, nil
}
