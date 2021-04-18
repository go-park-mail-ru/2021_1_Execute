package postgre_user_repository

import (
	"2021_1_Execute/internal/files"
	"2021_1_Execute/internal/users"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreUserRepository struct {
	Pool     *pgxpool.Pool
	logger   pgx.Logger
	FileUtil files.FileUtil
}

func NewPostgreUserRepository(pool *pgxpool.Pool, fileUtil files.FileUtil) users.UserRepository {
	return &PostgreUserRepository{
		Pool:     pool,
		logger:   pool.Config().ConnConfig.Logger,
		FileUtil: fileUtil,
	}
}

func (repo *PostgreUserRepository) log(ctx context.Context, logLevel int, msg, method string, data map[string]interface{}, err error) {
	repo.logger.Log(ctx, pgx.LogLevel(logLevel), msg, map[string]interface{}{
		"package": "postgre_user_repository",
		"method":  "PostgreUserRepository." + method,
		"data":    data,
		"error":   err,
	})
}
