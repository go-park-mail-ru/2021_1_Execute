package postgreRepo

import (
	"2021_1_Execute/src/api"

	"github.com/pkg/errors"
)

func (repo *PostgreRepo) IsAuthorized(session string) (api.User, bool, error) {
	//TODO: validation of session

	rows, err := repo.pool.Query("select user_id from sessions where session_token = $1::text", session)
	if err != nil {
		return api.User{}, false, errors.Wrap(err, "Unable to query authorization request")
	}

	var userID int = -1

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return api.User{}, false, errors.Wrap(err, "Unable to get user_id")
		}
	}

	rows.Close()

	if userID == -1 {
		return api.User{}, false, nil
	}

	user, err := repo.getUserByEmailOrID("ID", userID)
	if err != nil {
		return api.User{}, false, errors.Wrap(err, "Unable to get user")
	}

	return user, true, nil
}

func (repo *PostgreRepo) SetSession(session string, userID int) error {
	//TODO: validate session, userID
	rows, err := repo.pool.Query("insert into sessions (session_token, user_id) values ($1::text, $2::int)", session, userID)
	if err != nil {
		return errors.Wrap(err, "Unable to set session")
	}
	rows.Close()
	return nil
}

func (repo *PostgreRepo) DeleteSession(session string) error {
	//TODO: validate session
	rows, err := repo.pool.Query("delete from sessions where session_token = $1::text", session)
	if err != nil {
		return errors.Wrap(err, "Unable to delete session")
	}
	rows.Close()
	return nil
}
