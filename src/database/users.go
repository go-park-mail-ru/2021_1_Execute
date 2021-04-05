package postgreRepo

import (
	"2021_1_Execute/src/api"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (repo *PostgreRepo) getUserByEmailOrID(typeOfSelecting string, param interface{}) (api.User, error) {
	conn, err := repo.GetConnection()
	if err != nil {
		return api.User{}, err
	}
	defer conn.Close()

	var rows *pgx.Rows

	switch typeOfSelecting {
	case "email":
		rows, err = conn.Query("select id, email, username, hashed_password, path_to_avatar from users where email = $1", param.(string))
	case "ID":
		rows, err = conn.Query("select id, email, username, hashed_password, path_to_avatar from users where id = $1", param.(int))
	default:
		return api.User{}, errors.New("Invalid get request")
	}

	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while query getUserByEmailOrID")
	}
	defer rows.Close()

	var result api.User

	for rows.Next() {
		user, err := rows.Values()
		if err != nil {
			return api.User{}, errors.Wrap(err, "Error while reading getUserByEmailOrID")
		}

		if len(user) == 5 {
			result = api.User{
				ID:       int(user[0].(int32)),
				Email:    user[1].(string),
				Username: user[2].(string),
				Password: user[3].(string),
				Avatar:   user[4].(string),
			}
		}
	}
	return result, nil
}

func (repo *PostgreRepo) insertUser(user api.User) error {
	conn, err := repo.GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Query("insert into users (email, username, hashed_password, path_to_avatar) values ($1, $2, $3, $4)",
		user.Email,
		user.Username,
		user.Password,
		user.Avatar,
	)

	if err != nil {
		return errors.Wrap(err, "Error while query insertUser")
	}

	return nil
}

func (repo *PostgreRepo) CreateUser(input *api.UserRegistrationRequest) (api.User, error) {
	//TODO: add validation of username and path to avatar

	if !api.IsEmailValid(input.Email) {
		return api.User{}, api.BadRequestError
	}

	if !api.IsPasswordValid((*input).Password) {
		return api.User{}, api.BadRequestError
	}

	existingUser, err := repo.getUserByEmailOrID("email", input.Email)

	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while getUserByEmail in CreateUser")
	}
	if existingUser.Email != "" {
		return api.User{}, api.ConflictError
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return api.User{}, api.InternalServerError
	}

	user := api.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(passwordHashBytes),
	}

	err = repo.insertUser(user)

	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while inserting user in CreateUser")
	}

	return user, nil
}

func (repo *PostgreRepo) updateUserQuery(user api.User) error {
	conn, err := repo.GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (repo *PostgreRepo) UpdateUser(userID int, username, email, password, avatar string) error {
	switch {
	case email != "" && !api.IsEmailValid(email):
		return api.BadRequestError
	case email != "":
		ok, err := repo.IsEmailUniq(userID, email)
		if err != nil {
			return errors.Wrap(err, "Error while updating user")
		}
		if !ok {
			return api.ConflictError
		}
	case password != "" && !api.IsPasswordValid(password):
		return api.BadRequestError
	}

	userForUpdating, err := repo.getUserByEmailOrID("ID", userID)
	if err != nil {
		return errors.Wrap(err, "Error while updating user")
	}
	if userForUpdating.Email == "" {
		return api.NotFoundError
	}

	var user api.User

	user.ID = userID

	if username != "" {
		user.Username = username
	} else {
		user.Username = userForUpdating.Username
	}

	if email != "" {
		user.Email = email
	} else {
		user.Email = userForUpdating.Email
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			return api.InternalServerError
		}
		user.Password = string(hash)
	} else {
		user.Password = userForUpdating.Password
	}

	if avatar != "" {
		if userForUpdating.Avatar != "" {
			err := api.DeleteFile(userForUpdating.Avatar)
			if err != nil {
				return api.InternalServerError
			}
		}
		user.Avatar = avatar
	} else {
		user.Avatar = userForUpdating.Avatar
	}

	err = repo.updateUserQuery(user)

	if err != nil {
		return errors.Wrap(err, "Error while updating user")
	}

	return nil
}
