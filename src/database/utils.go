package postgreRepo

import (
	"2021_1_Execute/src/api"

	"github.com/pkg/errors"
)

func (repo *PostgreRepo) IsEmailUniq(userID int, email string) (bool, error) {
	if !api.IsEmailValid(email) {
		return false, api.BadRequestError
	}

	user, err := repo.getUserByEmailOrID("email", email)

	if err != nil {
		return false, errors.Wrap(err, "Error while checking uniq email")
	}

	if user.Email != "" {
		return false, nil
	}

	return true, nil
}

func (repo *PostgreRepo) getIdByEmail(email string) (int, error) {
	if !api.IsEmailValid(email) {
		return -1, api.BadRequestError
	}

	conn, err := repo.GetConnection()
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	rows, err := conn.Query("select id from users where email = $1::text", email)

	if err != nil {
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}
