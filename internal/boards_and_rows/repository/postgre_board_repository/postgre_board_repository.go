package postgre_board_repository

import (
	"2021_1_Execute/internal/boards_and_rows"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreBoardRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgreBoardRepository(pool *pgxpool.Pool) boards_and_rows.BoardRepository {
	return &PostgreBoardRepository{Pool: pool}
}
