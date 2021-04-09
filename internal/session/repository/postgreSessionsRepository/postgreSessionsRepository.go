package postgreSessionsRepository

import (
	"2021_1_Execute/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type PostgreSessionsRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgreSessionsRepository(pool *pgxpool.Pool) domain.SessionsRepository {
	return &PostgreSessionsRepository{Pool: pool}
}

func (repo *PostgreSessionsRepository) IsAuthorized(ctx context.Context, uuid string) (bool, int, error) {
	rows, err := repo.Pool.Query(ctx, "select user_id from sessions where session_token = $1::text", uuid)
	if err != nil {
		return false, -1, errors.Wrap(err, "Unable to query authorization request")
	}

	var userID int = -1

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return false, -1, errors.Wrap(err, "Unable to get user_id")
		}
	}

	rows.Close()

	if userID == -1 {
		return false, -1, nil
	}

	return true, userID, nil
}

func (repo *PostgreSessionsRepository) SetSession(ctx context.Context, userID int, uuid string) error {
	fmt.Println("hi")
	rows, err := repo.Pool.Query(ctx, "insert into sessions (session_token, user_id) values ($1::text, $2::int)", uuid, userID)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "Unable to set session")
	}
	rows.Close()
	return nil
}

func (repo *PostgreSessionsRepository) DeleteSession(ctx context.Context, uuid string) error {
	rows, err := repo.Pool.Query(ctx, "delete from sessions where session_token = $1::text", uuid)
	if err != nil {
		return errors.Wrap(err, "Unable to delete session")
	}
	rows.Close()
	return nil
}
