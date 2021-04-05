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

	conn, err := repo.GetConnection()
	if err != nil {
		return api.User{}, err
	}
	defer conn.Close()

	var rows *pgx.Rows

	switch typeOfSelecting {
	case "email":
		rows, err = conn.Query("select id, email, username, hashed_password, path_to_avatar from users where email = $1::text", param.(string))
	case "ID":
		rows, err = conn.Query("select id, email, username, hashed_password, path_to_avatar from users where id = $1::int", param.(int))
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
	conn, err := repo.GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Query("insert into users (email, username, hashed_password, path_to_avatar) values ($1::text, $2::text, $3::text, $4::text)",
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

	user.ID, err = repo.getIdByEmail(user.Email)

	if err != nil {
		return api.User{}, errors.Wrap(err, "Error while getting ID")
	}

	return user, nil
}

func (repo *PostgreRepo) updateUserQuery(user api.User) error {
	conn, err := repo.GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Query("update users set email = $1::text, username = $2::text, hashed_password = $3::text, path_to_avatar = $4::text where id = $5::int",
		user.Email,
		user.Username,
		user.Password,
		user.Avatar,
		user.ID,
	)

	if err != nil {
		return errors.Wrap(err, "Unable to update user")
	}

	return nil
}

func (repo *PostgreRepo) UpdateUser(userID int, username, email, password, avatar string) error {
	switch {
	//TODO: validation of username and path to avatar
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
