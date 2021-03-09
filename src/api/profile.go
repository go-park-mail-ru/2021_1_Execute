package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type GetUserByIdResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PatchUserRequest struct {
	NewEmail    string `json:"email,omitempty"`
	NewUsername string `json:"username,omitempty"`
	NewPassword string `json:"password,omitempty"`
}

func (db *Database) IsEmailUniq(userID int, email string) bool {
	for _, user := range *db.Users {
		if userID != user.ID && user.Email == email {
			return false
		}
	}
	return true
}

func (db *Database) UpdateUser(userID int, username, email, password string) error {
	switch {
	case email != "" && !IsEmailValid(email):
		return errors.New("Invalid email")
	case email != "" && !db.IsEmailUniq(userID, email):
		return errors.New("Non-uniq email")
	case password != "" && !IsPasswordValid(password):
		return errors.New("Invalid password")
	}

	for i, user := range *db.Users {
		if userID == user.ID {
			if username != "" {
				(*db.Users)[i].Username = username
			}
			if email != "" {
				(*db.Users)[i].Email = email
			}
			if password != "" {
				hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
				if err != nil {
					return errors.Wrap(err, "Error while hashing")
				}
				(*db.Users)[i].Password = string(hash)
				return nil
			}
		}
	}

	return errors.New("No such user")
}

func (db *Database) DeleteUser(userID int) error {
	for i, user := range *db.Users {
		if user.ID == userID {
			*db.Users = append((*db.Users)[:i], (*db.Users)[i+1:]...)
			return nil
		}
	}
	return errors.New("No such user")
}

func createGetUserByIdResponse(user User) GetUserByIdResponse {
	return GetUserByIdResponse{
		Email:    user.Email,
		Username: user.Username,
	}
}

func GetUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	ok, user := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
	}

	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, struct{ user GetUserByIdResponse }{createGetUserByIdResponse(user)})
}

func PatchUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	input := new(PatchUserRequest)

	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	ok, user := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err = db.UpdateUser(userID, input.NewUsername, input.NewEmail, input.NewPassword)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Invalid format").Error())
	}

	return c.NoContent(http.StatusOK)
}

func DeleteUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db := c.(*Database)

	ok, user := db.IsAuthorized(c)

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err = db.DeleteUser(userID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = DeleteSession(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
