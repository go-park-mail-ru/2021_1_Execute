package http

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/users/models"
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (handler *UserHandler) Login(c echo.Context) error {
	input := new(models.UserLoginRequest)
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	ctx := context.Background()
	userID, err := handler.userUC.Authentication(ctx, models.LoginRequestToUser(input))
	if err != nil {
		return err
	}

	err = handler.sessionHD.SetSession(c, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &models.EntranceResponse{ID: userID})
}

func (handler *UserHandler) Registration(c echo.Context) error {
	input := new(models.UserRegistrationRequest)
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	ctx := context.Background()
	userID, err := handler.userUC.Registration(ctx, models.RegistrationRequestToUser(input))
	if err != nil {
		return err
	}

	err = handler.sessionHD.SetSession(c, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &models.EntranceResponse{ID: userID})
}

func (handler *UserHandler) Logout(c echo.Context) error {
	err := handler.sessionHD.DeleteSession(c)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
