package postgre_board_repository

import (
	"2021_1_Execute/internal/boards_and_rows"
	"2021_1_Execute/internal/domain"
	"context"

	"github.com/pkg/errors"
)

func (repo *PostgreBoardRepository) AddBoard(ctx context.Context, board boards_and_rows.Board) (int, error) {
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

func (repo *PostgreBoardRepository) UpdateBoard(ctx context.Context, board boards_and_rows.Board) error {
	outdatedBoard, err := repo.GetBoard(ctx, board.ID)

	if err != nil {
		return errors.Wrap(err, "Unable to get outdated board")
	}

	newBoard := createUpdateBoardObject(outdatedBoard, board)

	err = repo.updateBoardQuery(ctx, newBoard)

	if err != nil {
		return errors.Wrap(err, "Unable to query updating request")
	}

	return nil
}

func createUpdateBoardObject(outdatedBoard, newBoard boards_and_rows.Board) boards_and_rows.Board {
	var result boards_and_rows.Board

	result.ID = outdatedBoard.ID

	if newBoard.Name == "" {
		result.Name = outdatedBoard.Name
	} else {
		result.Name = newBoard.Name
	}

	if newBoard.Description == "" {
		result.Description = outdatedBoard.Description
	} else {
		result.Description = newBoard.Description
	}

	return result
}

func (repo *PostgreBoardRepository) updateBoardQuery(ctx context.Context, board boards_and_rows.Board) error {
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

func (repo *PostgreBoardRepository) GetBoard(ctx context.Context, boardID int) (boards_and_rows.Board, error) {
	rows, err := repo.Pool.Query(ctx, "select id, name, description from boards where id = $1::int", boardID)

	if err != nil {
		return boards_and_rows.Board{}, errors.Wrap(err, "Unable to get board")
	}

	var board boards_and_rows.Board

	for rows.Next() {
		err = rows.Scan(&board.ID, &board.Name, &board.Description)
		if err != nil {
			return boards_and_rows.Board{}, errors.Wrap(err, "Unable to read board")
		}
	}

	rows.Close()

	if board.ID == 0 {
		return boards_and_rows.Board{}, domain.DBNotFoundError
	}

	return board, nil
}

func (repo *PostgreBoardRepository) GetUsersBoards(ctx context.Context, userID int) ([]boards_and_rows.Board, error) {
	rows, err := repo.Pool.Query(ctx,
		`select boards.id, boards.name, boards.description
		from owners
		where owners.user_id = $1::int and owners.board_id = boards.id
		left join administrators
		on administrators.user_id = $1::int and administrators.board_id = boards.id
		inner join boards`, userID) //todo is it works?

	if err != nil {
		return []boards_and_rows.Board{}, errors.Wrap(err, "Unable to get user's boards")
	}

	var boards []boards_and_rows.Board

	for rows.Next() {
		var board boards_and_rows.Board
		err = rows.Scan(&board.ID, &board.Name, &board.Description)

		if err != nil {
			return []boards_and_rows.Board{}, errors.Wrap(err, "Unable to get board")
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

func (repo *PostgreBoardRepository) GetBoardsOwner(ctx context.Context, boardID int) (int, error) {
	rows, err := repo.Pool.Query(ctx, "select user_id from owners where board_id = $1::int", boardID)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to get owner")
	}

	var owner int

	for rows.Next() {
		err = rows.Scan(&owner)
		if err != nil {
			return -1, errors.Wrap(err, "Unable to read owner")
		}
	}

	rows.Close()

	return owner, nil
}

func (repo *PostgreBoardRepository) GetBoardsAdmins(ctx context.Context, boardID int) ([]int, error) {
	rows, err := repo.Pool.Query(ctx, "select user_id from administrators where board_id = $1::int", boardID)
	if err != nil {
		return []int{}, errors.Wrap(err, "Unable to get admins")
	}

	var admins []int

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return []int{}, errors.Wrap(err, "Unable to read admin")
		}
		admins = append(admins, id)
	}

	rows.Close()

	return admins, nil
}

func (repo *PostgreBoardRepository) AddAdminToBoard(ctx context.Context, boardID int, userID int) error {
	rows, err := repo.Pool.Query(ctx, "insert into administrators (user_id, board_id) values ($1::int, $2::int)", userID, boardID)

	if err != nil {
		return errors.Wrap(err, "Unable to link user and board")
	}

	rows.Close()

	return nil
}

func (repo *PostgreBoardRepository) DeleteAdminFromBoard(ctx context.Context, boardID int, userID int) error {
	rows, err := repo.Pool.Query(ctx, "delete from administrators where user_id = $1::int and board_id = $2::int", userID, boardID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete admin")
	}
	rows.Close()

	return nil
}
