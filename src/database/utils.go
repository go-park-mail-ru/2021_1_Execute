package postgreRepo

import "github.com/pkg/errors"

func (repo *PostgreRepo) IsEmailUniq(userID int, email string) (bool, error) {
	user, err := repo.getUserByEmail(email)

	if err != nil {
		return false, errors.Wrap(err, "Error while checking uniq email")
	}

	if user.Email != "" {
		return false, nil
	}

	return true, nil
}
