package postgreUserRepository

import (
	"2021_1_Execute/internal/domain"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type PostgreUserRepository struct {
	Pool     *pgxpool.Pool
	FileUtil domain.FileUtil
}

func NewPostgreUserRepository(pool *pgxpool.Pool, fileUtil domain.FileUtil) domain.UserRepository {
	return &PostgreUserRepository{
		Pool:     pool,
		FileUtil: fileUtil,
	}
}

func (repo *PostgreUserRepository) AddUser(ctx context.Context, user domain.User) (int, error) {
	err := repo.insertUser(ctx, user)
	if err != nil {
		return -1, errors.Wrap(err, "Error while inserting user")
	}

	user.ID, err = repo.getIdByEmail(ctx, user.Email)
	if err != nil {
		return -1, errors.Wrap(err, "Error while getting ID")
	}

	return user.ID, nil
}

func (repo *PostgreUserRepository) UpdateUser(ctx context.Context, user domain.User) error {
	outdatedUser, err := repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return errors.Wrap(err, "Error while updating user")
	}
	if outdatedUser.Email == "" {
		return domain.NotFoundError
	}

	newUser, err := repo.createUserUpdateObject(outdatedUser, user)
	if err != nil {
		return errors.Wrap(err, "Unable to update user")
	}

	err = repo.updateUserQuery(ctx, newUser)
	if err != nil {
		return errors.Wrap(err, "Error while updating user")
	}

	return nil
}

func (repo *PostgreUserRepository) DeleteUser(ctx context.Context, userID int) error {
	user, err := repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "Unable to get user")
	}
	if user.Email == "" {
		return domain.NotFoundError
	}

	rows, err := repo.Pool.Query(ctx, "delete from users where id = $1::int", user.ID)
	if err != nil {
		return errors.Wrap(err, "Unable to delete user")
	}
	rows.Close()

	return nil
}

func (repo *PostgreUserRepository) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	rows, err := repo.Pool.Query(ctx, "select id, email, username, hashed_password, path_to_avatar from users where id = $1::int", userID)

	if err != nil {
		return domain.User{}, errors.Wrap(err, "Error while getting user by id")
	}
	defer rows.Close()

	var result domain.User

	for rows.Next() {
		err := rows.Scan(&result.ID, &result.Email, &result.Username, &result.Password, &result.Avatar)
		if err != nil {
			return domain.User{}, errors.Wrap(err, "Error while reading user")
		}
	}
	return result, nil
}

func (repo *PostgreUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	rows, err := repo.Pool.Query(ctx, "select id, email, username, hashed_password, path_to_avatar from users where email = $1::text", email)

	if err != nil {
		return domain.User{}, errors.Wrap(err, "Error while getting user by email")
	}
	defer rows.Close()

	var result domain.User

	for rows.Next() {
		err := rows.Scan(&result.ID, &result.Email, &result.Username, &result.Password, &result.Avatar)
		if err != nil {
			return domain.User{}, errors.Wrap(err, "Error while reading user")
		}
	}
	return result, nil
}

func (repo *PostgreUserRepository) insertUser(ctx context.Context, user domain.User) error {
	rows, err := repo.Pool.Query(ctx, "insert into users (email, username, hashed_password, path_to_avatar) values ($1::text, $2::text, $3::text, $4::text)",
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

func (repo *PostgreUserRepository) getIdByEmail(ctx context.Context, email string) (int, error) {
	rows, err := repo.Pool.Query(ctx, "select id from users where email = $1::text", email)
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

func (repo *PostgreUserRepository) updateUserQuery(ctx context.Context, user domain.User) error {
	rows, err := repo.Pool.Query(ctx, "update users set email = $1::text, username = $2::text, hashed_password = $3::text, path_to_avatar = $4::text where id = $5::int",
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

func (repo *PostgreUserRepository) createUserUpdateObject(outdatedUser, newUser domain.User) (domain.User, error) {
	var result domain.User

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
		result.Password = newUser.Password
	} else {
		result.Password = outdatedUser.Password
	}

	if newUser.Avatar != "" {
		if outdatedUser.Avatar != "" {
			err := repo.FileUtil.DeleteFile(outdatedUser.Avatar)
			if err != nil {
				return domain.User{}, err
			}
		}
		result.Avatar = newUser.Avatar
	} else {
		result.Avatar = outdatedUser.Avatar
	}

	return result, nil
}
