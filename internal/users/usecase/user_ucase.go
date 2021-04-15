package usecase

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/users"
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo users.UserRepository
}

func NewUserUsecase(repo users.UserRepository) users.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func setPassword(user users.User) (users.User, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return users.User{}, errors.Wrap(domain.InternalServerError, "Error while hashing:"+err.Error())
	}

	user.Password = string(passwordHashBytes)
	return user, nil
}

func (uc *userUsecase) Registration(ctx context.Context, user users.User) (int, error) {
	user, err := setPassword(user)
	userId, err := uc.userRepo.AddUser(ctx, user)
	if err != nil {
		return 0, domain.DBErrorToServerError(err)
	}
	return userId, nil
}

func (uc *userUsecase) UpdateAvatar(ctx context.Context, userID int, path string) error {
	changedUser := users.User{
		ID:     userID,
		Avatar: path,
	}
	err := uc.userRepo.UpdateUser(ctx, changedUser)
	return domain.DBErrorToServerError(err)
}

func (uc *userUsecase) UpdateUser(ctx context.Context, changerID int, changedUser users.User) error {
	var err error
	if changerID != changedUser.ID {
		return errors.Wrap(domain.ForbiddenError, "Not enough rights")
	}

	if len(changedUser.Password) > 0 {
		changedUser, err = setPassword(changedUser)
		if err != nil {
			return err
		}
	}
	err = uc.userRepo.UpdateUser(ctx, changedUser)
	return domain.DBErrorToServerError(err)

}

func (uc *userUsecase) DeleteUser(ctx context.Context, changerID int, userID int) error {
	if changerID != userID {
		return errors.Wrap(domain.ForbiddenError, "Not enough rights")
	}
	err := uc.userRepo.DeleteUser(ctx, userID)

	return domain.DBErrorToServerError(err)
}

func (uc *userUsecase) GetUserByID(ctx context.Context, userID int) (users.User, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	return user, domain.DBErrorToServerError(err)
}

func (uc *userUsecase) Authentication(ctx context.Context, user users.User) (int, error) {
	userFromBD, err := uc.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, domain.DBErrorToServerError(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(userFromBD.Password), []byte(user.Password)) == nil {

		return userFromBD.ID, nil
	}
	return 0, errors.Wrap(domain.ForbiddenError, "Wrong pair: email, password")
}
