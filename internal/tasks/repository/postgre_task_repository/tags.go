package postgre_task_repository

import (
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreTaskRepository) AddTag(ctx context.Context, taskID int, tag tasks.Tag) (int, error) {
	rows, err := repo.Pool.Query(ctx, "insert into task (color, name) values ($1::text, $2::text) returning id", tag.Color, tag.Name)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to insert tag")
	}
	defer rows.Close()

	var tagID int = -1

	for rows.Next() {
		err = rows.Scan(&tagID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read tag id")
		}
	}

	if tagID == -1 {
		return -1, errors.New("Invalid tag id")
	}

	return tagID, nil
}

func (repo *PostgreTaskRepository) DeleteTag(ctx context.Context, tagID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from tags where id = $1::int", tagID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete tag")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) AddTagToTask(ctx context.Context, taskID, tagID int) error {
	rows, err := repo.Pool.Query(ctx, "insert into tags_tasks (tag_id, task_id) values ($1::int, $2::int)", tagID, taskID)
	if err != nil {
		return errors.Wrap(err, "Unable to ling tag and task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) AddTagToBoard(ctx context.Context, boardID, tagID int) error {
	rows, err := repo.Pool.Query(ctx, "update tags set board_id = $1::int where id = $2::int", boardID, tagID)
	if err != nil {
		return errors.Wrap(err, "Unable to link tag and board")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) DeleteTagFromTask(ctx context.Context, taskID, tagID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from tags_tasks where tag_id = $1::int and task_id = $2::int", tagID, taskID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete tag from task")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) DeleteTagFromBoard(ctx context.Context, tagID int) error {
	rows, err := repo.Pool.Query(ctx, "update tags set board_id = null where id = $1::int", tagID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete tag from board")
	}

	rows.Close()

	return nil
}

func (repo *PostgreTaskRepository) GetTasksTags(ctx context.Context, taskID int) ([]tasks.Tag, error) {
	rows, err := repo.Pool.Query(ctx, `select (tags.id, tags.name, tags.color) from tags
	inner join tags_tasks as tt
	on tt.task_id = $1::int and tt.tag_id = tags.id`, taskID)
	if err != nil {
		return []tasks.Tag{}, errors.Wrap(err, "Unable to get task's tags")
	}
	defer rows.Close()

	var tags []tasks.Tag

	for rows.Next() {
		var tag tasks.Tag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Color)
		if err != nil {
			return []tasks.Tag{}, errors.Wrap(err, "Unable to read task's tag")
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (repo *PostgreTaskRepository) GetBoardsTags(ctx context.Context, boardID int) ([]tasks.Tag, error) {
	rows, err := repo.Pool.Query(ctx, "select (id, name, color) from tags where board_id = $1::int", boardID)
	if err != nil {
		return []tasks.Tag{}, errors.Wrap(err, "Unable to get board's tags")
	}
	defer rows.Close()

	var tags []tasks.Tag

	for rows.Next() {
		var tag tasks.Tag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Color)
		if err != nil {
			return []tasks.Tag{}, errors.Wrap(err, "Unable to read board's tag")
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
