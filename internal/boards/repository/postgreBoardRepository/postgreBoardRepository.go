package postgreBoardRepository

import (
	"2021_1_Execute/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreBoardRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgreBoardRepository(pool *pgxpool.Pool) domain.BoardRepository {
	return &PostgreBoardRepository{Pool: pool}
}
