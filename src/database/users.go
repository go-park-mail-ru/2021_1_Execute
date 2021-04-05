package postgreRepo

import (
	"2021_1_Execute/src/api"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (repo *PostgreRepo) getUserByEmailOrID(typeOfSelecting string, param interface{}) (api.User, error) {
	if typeOfSelecting == "email" && !api.IsEmailValid(param.(string)) {
		return api.User{}, api.BadRequestError
	}

	var rows *pgx.Rows
	var err error

	switch typeOfSelecting {
	case "email":
		rows, err = repo.pool.Query("select id, email, username, hashed_password, path_to_avatar from users where email = $1::text", param.(string))
	case "ID":
		rows, err = repo.pool.Query("select id, email, username, hashed_password, path_to_avatar from users where id = $1::int", param.(int))
	default:
		return api.User{}, errors.New("Invalid get request")
	}

	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while query getUserByEmailOrID")
	}
	defer rows.Close()

	var result api.User

	for rows.Next() {
		err := rows.Scan(&result.ID, &result.Email, &result.Username, &result.Password, &result.Avatar)
		if err != nil {
			return api.User{}, errors.Wrap(err, "Error while reading getUserByEmailOrID")
		}
	}
	return result, nil
}

func (repo *PostgreRepo) insertUser(user api.User) error {
	rows, err := repo.pool.Query("insert into users (email, username, hashed_password, path_to_avatar) values ($1::text, $2::text, $3::text, $4::text)",
		user.Email,
		user.Username,
		user.Password,
		user.Avatar,
	)

	if err != nil {
		return errors.Wrap(err, "Error while query insertUser")
	}
	rows.Close()

	return nil
}

func createUserInsertObject(input *api.UserRegistrationRequest) (api.User, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return api.User{}, errors.Wrap(err, "Unable to insert user")
	}

	return api.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(passwordHashBytes),
	}, nil
}

func (repo *PostgreRepo) CreateUser(input *api.UserRegistrationRequest) (api.User, error) {
	//TODO: add validation of username and path to avatar
	if !api.IsEmailValid(input.Email) {
		return api.User{}, api.BadRequestError
	}
	if !api.IsPasswordValid((*input).Password) {
		return api.User{}, api.BadRequestError
	}

	uniq, err := repo.IsEmailUniq(input.Email)
	if err != nil {
		return api.User{}, errors.Wrap(err, "Unable to check uniqueness")
	}
	if !uniq {
		return api.User{}, api.ConflictError
	}

	user, err := createUserInsertObject(input)
	if err != nil {
		return api.User{}, errors.Wrap(err, "Unable to create user")
	}

	err = repo.insertUser(user)
	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while inserting user in CreateUser")
	}

	user.ID, err = repo.getIdByEmail(user.Email)
	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while getting ID")
	}

	return user, nil
}

func (repo *PostgreRepo) updateUserQuery(user api.User) error {
	rows, err := repo.pool.Query("update users set email = $1::text, username = $2::text, hashed_password = $3::text, path_to_avatar = $4::text where id = $5::int",
		user.Email,
		user.Username,
		user.Password,
		user.Avatar,
		user.ID,
	)

	if err != nil {
		return errors.Wrap(err, "Unable to update user")
	}
	rows.Close()

	return nil
}

func createUserUpdateObject(outdatedUser, newUser api.User) (api.User, error) {
	var result api.User

	result.ID = outdatedUser.ID
	if newUser.Username != "" {
		result.Username = newUser.Username
	} else {
		result.Username = outdatedUser.Username
	}

	if newUser.Email != "" {
		result.Email = newUser.Email
	} else {
		result.Email = outdatedUser.Email
	}

	if newUser.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
		if err != nil {
			return api.User{}, err
		}
		result.Password = string(hash)
	} else {
		result.Password = outdatedUser.Password
	}

	if newUser.Avatar != "" {
		if outdatedUser.Avatar != "" {
			err := api.DeleteFile(outdatedUser.Avatar)
			if err != nil {
				return api.User{}, err
			}
		}
		result.Avatar = newUser.Avatar
	} else {
		result.Avatar = outdatedUser.Avatar
	}

	return result, nil
}

func (repo *PostgreRepo) UpdateUser(userID int, username, email, password, avatar string) error {
	switch {
	//TODO: validation of username and path to avatar
	case email != "" && !api.IsEmailValid(email):
		return api.BadRequestError
	case email != "":
		ok, err := repo.IsEmailUniq(email)
		if err != nil {
			return errors.Wrap(err, "Error while updating user")
		}
		if !ok {
			return api.ConflictError
		}
	case password != "" && !api.IsPasswordValid(password):
		return api.BadRequestError
	}

	outdatedUser, err := repo.getUserByEmailOrID("ID", userID)
	if err != nil {
		return errors.Wrap(err, "Error while updating user")
	}
	if outdatedUser.Email == "" {
		return api.NotFoundError
	}

	newUser, err := createUserUpdateObject(outdatedUser, api.User{
		Username: username,
		Email:    email,
		Password: password,
		Avatar:   avatar,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to update user")
	}

	err = repo.updateUserQuery(newUser)
	if err != nil {
		return errors.Wrap(err, "Error while updating user")
	}

	return nil
}

func (repo *PostgreRepo) DeleteUser(userID int) error {
	//TODO: id validation
	user, err := repo.getUserByEmailOrID("ID", userID)
	if err != nil {
		return errors.Wrap(err, "Unable to get user")
	}
	if user.Email == "" {
		return api.NotFoundError
	}

	rows, err := repo.pool.Query("delete from users where id = $1::int", userID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete user")
	}
	rows.Close()

	return nil
}
