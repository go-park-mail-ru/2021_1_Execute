package http

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

var Ls string = "dfs"

type UserLoginRequest struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" ` // todo valid:"password"
}

type EntranceResponse struct {
	ID int `json:"id"`
}

type UserRegistrationRequest struct {
	Email    string `json:"email" valid:"email"`
	Username string `json:"username" ` //todo valid:"username"
	Password string `json:"password" ` // todovalid:"password"
}

func LoginRequestToUser(input *UserLoginRequest) domain.User {
	return domain.User{
		Email:    input.Email,
		Password: input.Password,
	}
}

func RegistrationRequestToUser(input *UserRegistrationRequest) domain.User {
	return domain.User{
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	}
}

func (handler *UserHandler) Login(c echo.Context) error {
	input := new(UserLoginRequest)
	if err := c.Bind(input); err != nil {
		return domain.GetEchoError(errors.Wrap(domain.BadRequestError, err.Error()))
	}

	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return domain.GetEchoError(errors.Wrap(domain.BadRequestError, err.Error()))
	}

	ctx := context.Background()
	userID, err := handler.userUC.Authentication(ctx, LoginRequestToUser(input))
	if err != nil {
		return domain.GetEchoError(err)
	}

	err = handler.sessionHD.SetSession(c, userID)
	if err != nil {
		return domain.GetEchoError(err)
	}

	return c.JSON(http.StatusOK, &EntranceResponse{ID: userID})
}

func (handler *UserHandler) Registration(c echo.Context) error {
	input := new(UserRegistrationRequest)
	if err := c.Bind(input); err != nil {
		return domain.GetEchoError(errors.Wrap(domain.BadRequestError, err.Error()))
	}

	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return domain.GetEchoError(errors.Wrap(domain.BadRequestError, err.Error()))
	}

	ctx := context.Background()
	userID, err := handler.userUC.Registration(ctx, RegistrationRequestToUser(input))
	if err != nil {
		return domain.GetEchoError(err)
	}

	err = handler.sessionHD.SetSession(c, userID)
	if err != nil {
		return domain.GetEchoError(err)
	}
	return c.JSON(http.StatusOK, &EntranceResponse{ID: userID})
}

func (handler *UserHandler) Logout(c echo.Context) error {
	err := handler.sessionHD.DeleteSession(c)
	if err != nil {
		return domain.GetEchoError(err)
	}

	return c.NoContent(http.StatusOK)
}
