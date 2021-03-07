package api

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (db *Database) CreateUser(input *UserRegistrationRequest) (User, error, int) {
	if !IsEmailValid(input.Email) {
		return User{}, errors.New("Invalid email"), 400
	}
	if !IsPasswordValid((*input).Password) {
		return User{}, errors.New("Invalid password"), 400
	}
	for _, user := range *db.Users {
		if user.Email == input.Email {
			return User{}, errors.New("Email not unique"), 409
		}
	}
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, errors.Wrap(err, "Error while hashing"), 500
	}
	newUser := User{
		ID:       uint64(len(*db.Users)),
		Email:    input.Email,
		Username: input.Username,
		Password: string(passwordHashBytes),
	}
	*db.Users = append(*db.Users, newUser)
	return newUser, nil, 200
}

func (db *Database) IsCredentialsCorrect(input *UserLoginRequest) (bool, User) {
	for _, user := range *db.Users {
		if user.Email == input.Email && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) == nil {
			return true, user
		}
	}
	return false, User{}
}

func (db *Database) IsAuthorized(c echo.Context) (bool, User) {
	session, err := c.Cookie(CookieName)
	if err != nil {
		return false, User{}
	}
	userID, isAuthorized := (*db.Sessions)[session.Value]
	if isAuthorized {
		for _, user := range *db.Users {
			if user.ID == userID {
				return true, user
			}
		}
	}
	return false, User{}
}
