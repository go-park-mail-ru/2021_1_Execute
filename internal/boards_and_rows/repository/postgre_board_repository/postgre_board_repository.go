package postgre_board_repository

import (
	"2021_1_Execute/internal/boards_and_rows"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	logLevelTrace = 6
	logLevelDebug = 5
	logLevelInfo  = 4
	logLevelWarn  = 3
	logLevelError = 2
	logLevelNone  = 1
)

type PostgreBoardRepository struct {
	Pool   *pgxpool.Pool
	logger pgx.Logger
}

func NewPostgreBoardRepository(pool *pgxpool.Pool) boards_and_rows.BoardRepository {
	return &PostgreBoardRepository{Pool: pool, logger: pool.Config().ConnConfig.Logger}
}
