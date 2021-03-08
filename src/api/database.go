package api

import (
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func (db *Database) CreateUser(input *UserRegistrationRequest) (User, error) {
	if !IsEmailValid(input.Email) {
		return User{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}
	if !IsPasswordValid((*input).Password) {
		return User{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid password")
	}

	for _, user := range *db.Users {
		if user.Email == input.Email {
			return User{}, echo.NewHTTPError(http.StatusConflict, "Email not unique")
		}
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, echo.NewHTTPError(http.StatusInternalServerError, "Error while hashing: "+err.Error())
	}

	newUser := User{
		ID:       int(len(*db.Users)),
		Email:    input.Email,
		Username: input.Username,
		Password: string(passwordHashBytes),
	}

	*db.Users = append(*db.Users, newUser)
	return newUser, nil
}

func (db *Database) IsCredentialsCorrect(input *UserLoginRequest) (User, bool) {
	for _, user := range *db.Users {
		if user.Email == input.Email && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) == nil {
			return user, true
		}
	}
	return User{}, false
}

func (db *Database) IsAuthorized(c echo.Context) (User, bool) {
	session, err := c.Cookie(CookieName)
	if err != nil {
		return User{}, false
	}

	userID, isAuthorized := (*db.Sessions)[session.Value]
	if isAuthorized {
		for _, user := range *db.Users {
			if user.ID == userID {
				return user, true
			}
		}
	}
	return User{}, false
}
