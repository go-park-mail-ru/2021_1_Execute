package postgreBoardRepository

import (
	"2021_1_Execute/internal/domain"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreBoardRepository) AddBoard(ctx context.Context, board domain.Board) (int, error) {
	rows, err := repo.Pool.Query(ctx, "insert into boards (name, description) values ($1::text, $2::text) returning id", board.Name, board.Description)

	if err != nil {
		return -1, errors.Wrap(err, "Unable to insert board")
	}

	var boardID int = -1

	for rows.Next() {
		err = rows.Scan(&boardID)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to get board id")
		}
	}

	rows.Close()

	return boardID, nil
}

func (repo *PostgreBoardRepository) AddOwner(ctx context.Context, boardID int, userID int) error {
	rows, err := repo.Pool.Query(ctx, "insert into owners (user_id, board_id) values ($1::int, $2::int)", userID, boardID)

	if err != nil {
		return errors.Wrap(err, "Unable to link user and board")
	}

	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) UpdateBoard(ctx context.Context, board domain.Board) error {
	rows, err := repo.Pool.Query(ctx, "update boards set name = $1::text, description = $2::text where id = $3::int",
		board.Name,
		board.Description,
		board.ID,
	)

	if err != nil {
		return errors.Wrap(err, "Unable to update board")
	}

	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) GetBoard(ctx context.Context, boardID int) (domain.Board, error) {
	rows, err := repo.Pool.Query(ctx, "select id, name, description from boards where id = $1::int", boardID)

	if err != nil {
		return domain.Board{}, errors.Wrap(err, "Unable to get board")
	}

	var board domain.Board

	for rows.Next() {
		err = rows.Scan(&board.ID, &board.Name, &board.Description)
		if err != nil {
			return domain.Board{}, errors.Wrap(err, "Unable to read board")
		}
	}

	if board.ID == 0 {
		return domain.Board{}, domain.DBNotFoundError
	}

	return board, nil
}

func (repo *PostgreBoardRepository) GetUsersBoards(ctx context.Context, userID int) ([]domain.Board, error) {
	rows, err := repo.Pool.Query(ctx,
		`select boards.id, boards.name, boards.description
	from boards
	inner join owners
	on owners.user_id = $1::int and owners.board_id = boards.id`, userID)

	if err != nil {
		return []domain.Board{}, errors.Wrap(err, "Unable to get user's boards")
	}

	var boards []domain.Board

	for rows.Next() {
		var board domain.Board
		err = rows.Scan(&board.ID, &board.Name, &board.Description)

		if err != nil {
			return []domain.Board{}, errors.Wrap(err, "Unable to get board")
		}

		boards = append(boards, board)
	}

	rows.Close()

	return boards, nil
}

func (repo *PostgreBoardRepository) DeleteBoard(ctx context.Context, boardID int) error {
	board, err := repo.GetBoard(ctx, boardID)
	if err != nil {
		return err
	}

	rows, err := repo.Pool.Query(ctx, "delete from boards where id = $1::int", board.ID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete board")
	}
	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) GetBoardOwners(ctx context.Context, boardID int) ([]int, error) {
	rows, err := repo.Pool.Query(ctx, "select user_id from owners where board_id = $1::int", boardID)
	if err != nil {
		return []int{}, errors.Wrap(err, "Unable to get owners")
	}

	var owners []int

	for rows.Next() {
		var owner int
		err = rows.Scan(&owner)
		if err != nil {
			return []int{}, errors.Wrap(err, "Unable to read owner")
		}
		owners = append(owners, owner)
	}

	rows.Close()

	return owners, nil
}