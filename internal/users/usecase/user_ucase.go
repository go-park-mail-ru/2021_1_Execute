package usecase

import (
	"2021_1_Execute/internal/domain"
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func setPassword(user domain.User) (domain.User, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return domain.User{}, errors.Wrap(domain.InternalServerError, "Error while hashing:"+err.Error())
	}

	user.Password = string(passwordHashBytes)
	return user, nil
}

func (uc *userUsecase) Registration(ctx context.Context, user domain.User) (int, error) {
	user, err := setPassword(user)
	userId, err := uc.userRepo.AddUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (uc *userUsecase) UpdateUser(ctx context.Context, changerID int,
	changedUser domain.User, isPasswordChanged bool) error {
	var err error
	if changerID != changedUser.ID {
		return errors.Wrap(domain.ForbiddenError, "Not enough rights")
	}

	if isPasswordChanged {
		changedUser, err = setPassword(changedUser)
		if err != nil {
			return err
		}
	}
	err = uc.userRepo.UpdateUser(ctx, changedUser)
	return err

}

func (uc *userUsecase) DeleteUser(ctx context.Context, changerID int, userID int) error {
	if changerID != userID {
		return errors.Wrap(domain.ForbiddenError, "Not enough rights")
	}
	err := uc.userRepo.DeleteUser(ctx, userID)
	return err
}

func (uc *userUsecase) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	user, err := uc.GetUserByID(ctx, userID)
	return user, err
}

func (uc *userUsecase) Authentication(ctx context.Context, user domain.User) (int, error) {
	userFromBD, err := uc.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}
	if bcrypt.CompareHashAndPassword([]byte(userFromBD.Password), []byte(user.Password)) == nil {

		return userFromBD.ID, nil
	}
	return 0, errors.Wrap(domain.ForbiddenError, "Wrong pair: email, password")
}
