package api

import (
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
	for i := 0; i < len(*db.Users); i++ {
		if userID == (*db.Users)[i].ID {
			if input.NewEmail != "" {
				(*db.Users)[i].Email = input.NewEmail
			}
			if input.NewUsername != "" {
				(*db.Users)[i].Username = input.NewUsername
			}
			if input.NewPassword != "" {
				hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
				if err != nil {
					return nil
				}
				(*db.Users)[i].Password = string(hash)
			}
			break
		}
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
	for i, el := range *db.Users {
		if el.ID == userID {
			*db.Users = append((*db.Users)[:i], (*db.Users)[i+1:]...)
			err := DeteleSesssion(c)
			if err != nil {
				return err
			}
			break
		}
	}
	return c.NoContent(http.StatusOK)
}
