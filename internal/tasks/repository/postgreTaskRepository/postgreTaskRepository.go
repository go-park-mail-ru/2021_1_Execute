package postgreTaskRepository

import (
	"2021_1_Execute/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreTaskRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgreTaskRepository(pool *pgxpool.Pool) domain.TaskRepository {
	return &PostgreTaskRepository{Pool: pool}
}
