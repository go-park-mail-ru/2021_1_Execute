package postgreRepo

import (
	"2021_1_Execute/src/api"

	"github.com/pkg/errors"
)

func (repo *PostgreRepo) IsEmailUniq(userID int, email string) (bool, error) {
	if !api.IsEmailValid(email) {
		return false, api.BadRequestError
	}

	user, err := repo.getUserByEmail(email)

	if err != nil {
		return false, errors.Wrap(err, "Error while checking uniq email")
	}

	if user.Email != "" {
		return false, nil
	}

	return true, nil
}
