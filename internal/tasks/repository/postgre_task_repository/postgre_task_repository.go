package postgre_task_repository

import (
	"2021_1_Execute/internal/tasks"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreTaskRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgreTaskRepository(pool *pgxpool.Pool) tasks.TaskRepository {
	return &PostgreTaskRepository{Pool: pool}
}
