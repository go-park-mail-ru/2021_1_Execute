package postgreRepo

import (
	"2021_1_Execute/src/api"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (repo *PostgreRepo) IsEmailUniq(email string) (bool, error) {
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

	rows, err := repo.pool.Query("select id from users where email = $1::text", email)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (repo *PostgreRepo) IsCredentialsCorrect(input *api.UserLoginRequest) (api.User, bool, error) {
	user, err := repo.getUserByEmailOrID("email", input.Email)
	if err != nil {
		return api.User{}, false, errors.Wrap(err, "Unable to get user")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) == nil {
		return user, true, nil
	}
	return api.User{}, false, nil
}

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
