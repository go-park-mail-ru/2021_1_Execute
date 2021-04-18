package postgre_task_repository

import (
	"2021_1_Execute/internal/tasks"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreTaskRepository struct {
	Pool   *pgxpool.Pool
	logger pgx.Logger
}

func NewPostgreTaskRepository(pool *pgxpool.Pool) tasks.TaskRepository {
	return &PostgreTaskRepository{Pool: pool, logger: pool.Config().ConnConfig.Logger}
}

func (repo *PostgreTaskRepository) log(ctx context.Context, logLevel int, msg, method string, data map[string]interface{}, err error) {
	repo.logger.Log(ctx, pgx.LogLevel(logLevel), msg, map[string]interface{}{
		"package": "postgre_task_repository",
		"method":  "PostgreTaskRepository." + method,
		"data":    data,
		"error":   err,
	})
}
