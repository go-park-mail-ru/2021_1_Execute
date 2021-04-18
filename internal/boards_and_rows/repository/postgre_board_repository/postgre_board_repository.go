package postgre_board_repository

import (
	"2021_1_Execute/internal/boards_and_rows"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreBoardRepository struct {
	Pool   *pgxpool.Pool
	logger pgx.Logger
}

func NewPostgreBoardRepository(pool *pgxpool.Pool) boards_and_rows.BoardRepository {
	return &PostgreBoardRepository{Pool: pool, logger: pool.Config().ConnConfig.Logger}
}

func (repo *PostgreBoardRepository) log(ctx context.Context, logLevel int, msg, method string, data map[string]interface{}, err error) {
	repo.logger.Log(ctx, pgx.LogLevel(logLevel), msg, map[string]interface{}{
		"package": "postgre_board_repository",
		"method":  "PostgreBoardRepository." + method,
		"data":    data,
		"error":   err,
	})
}
