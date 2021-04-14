package postgre_task_repository

import (
	"2021_1_Execute/internal/tasks"

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

type PostgreTaskRepository struct {
	Pool   *pgxpool.Pool
	logger pgx.Logger
}

func NewPostgreTaskRepository(pool *pgxpool.Pool) tasks.TaskRepository {
	return &PostgreTaskRepository{Pool: pool, logger: pool.Config().ConnConfig.Logger}
}
