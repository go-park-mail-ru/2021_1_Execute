package http

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/session"
	"2021_1_Execute/internal/users"
	"2021_1_Execute/internal/users/models"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

<<<<<<< HEAD
type GetUserByIdResponse struct {
	Email     string `json:"email" `
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
}

type GetUserByIdBody struct {
	Response GetUserByIdResponse `json:"user"`
}

type PatchUserRequest struct {
	NewEmail    string `json:"email,omitempty" valid:"email"`
	NewUsername string `json:"username,omitempty" valid:"username"`
	NewPassword string `json:"password,omitempty" valid:"password"`
}

=======
>>>>>>> architecture
type UserHandler struct {
	userUC    users.UserUsecase
	sessionHD session.SessionHandler
}

func NewUserHandler(e *echo.Echo, userUsecase users.UserUsecase, sessionsHandler session.SessionHandler) {
	handler := &UserHandler{
		userUC:    userUsecase,
		sessionHD: sessionsHandler,
	}
	e.GET("/api/users/", handler.GetCurrentUser)
	e.GET("/api/users/:id", handler.GetUserByID)
	e.PATCH("/api/users/", handler.PatchUser)
	e.DELETE("/api/users/:id", handler.DeleteUserByID)
	e.POST("/api/login/", handler.Login)
	e.POST("/api/users/", handler.Registration)
	e.DELETE("/api/logout/", handler.Logout)
	e.GET("/api/authorized/", handler.IsAuthorized)
}

func (handler *UserHandler) IsAuthorized(c echo.Context) error {
	_, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (handler *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := handler.userUC.GetUserByID(ctx, userId)
	return c.JSON(http.StatusOK, models.CreateGetUserByIdResponse(user))
}

func (handler *UserHandler) GetUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return domain.IDFormatError
	}

	_, err = handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := handler.userUC.GetUserByID(ctx, userID)
	return c.JSON(http.StatusOK, models.CreateGetUserByIdResponse(user))
}

func (handler *UserHandler) PatchUser(c echo.Context) error {
	input := new(models.PatchUserRequest)
	if err := c.Bind(input); err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return errors.Wrap(domain.BadRequestError, err.Error())
	}

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user := models.CreateUserFromPatchRequest(input)
	user.ID = userID
	err = handler.userUC.UpdateUser(ctx, userID, user)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (handler *UserHandler) DeleteUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return domain.IDFormatError
	}

	currentUserID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = handler.userUC.DeleteUser(ctx, currentUserID, userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
