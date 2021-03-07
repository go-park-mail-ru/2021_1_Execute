package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type UserGetResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserPatchRequest struct {
	NewEmail    string `json:"email,omitempty"`
	NewUsername string `json:"username,omitempty"`
	NewPassword string `json:"password,omitempty"`
}

func (db *Database) UpdateUser(userID uint64, username, email, password string) error {
	for i := 0; i < len(*db.Users); i++ {
		if userID == (*db.Users)[i].ID {
			if username != "" {
				(*db.Users)[i].Username = username
			}
			if email != "" {
				(*db.Users)[i].Email = email
			}
			if password != "" {
				hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
				if err != nil {
					return err
				}
				(*db.Users)[i].Password = string(hash)
				return nil
			}
		}
	}
	return errors.New("No such user")
}

func (db *Database) DeleteUser(userID uint64) error {
	for i, user := range *db.Users {
		if user.ID == userID {
			*db.Users = append((*db.Users)[:i], (*db.Users)[i+1:]...)
			return nil
		}
	}
	return errors.New("No such user")
}

func createGetUserResponse(user User) UserGetResponse {
	return UserGetResponse{
		Email:    user.Email,
		Username: user.Username,
	}
}

func GetUserByID(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	db := c.(*Database)
	ok, user := db.IsAuthorized(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
	}
	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, createGetUserResponse(user))
}

func PatchUserByID(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	db := c.(*Database)
	ok, user := db.IsAuthorized(c)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}
	if userID != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	input := new(UserPatchRequest)
	if err := c.Bind(input); err != nil {
		return err
	}
	switch {
	case input.NewEmail != "" && !IsEmailValid(input.NewEmail):
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong format of email")
	case input.NewPassword != "" && !IsPasswordValid(input.NewPassword):
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong format of password")
	}
	err = db.UpdateUser(userID, input.NewUsername, input.NewEmail, input.NewPassword)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func DeleteUserByID(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
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
		return nil
	}
	err = DeteleSesssion(c)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
