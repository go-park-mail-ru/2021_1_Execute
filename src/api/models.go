package api

import (
	"regexp"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type SessionsMap map[string]uint64

type Database struct {
	echo.Context
	Users    *[]User
	Sessions *SessionsMap
}

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type RegistrationResponse struct {
	ID uint64 `json:"id"`
}

func CreateRegistrationResponse(user User) RegistrationResponse {
	return RegistrationResponse{ID: user.ID}
}

func (db *Database) CreateUser(input *UserRegistrationRequest) (User, error, int) {
	if !IsEmailValid((*input).Email) {
		return User{}, errors.New("Invalid email"), 400
	}
	if !IsPasswordValid((*input).Password) {
		return User{}, errors.New("Invalid password"), 400
	}
	for _, user := range *db.Users {
		if user.Email == (*input).Email {
			return User{}, errors.New("Email not unique"), 409
		}
	}
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, errors.Wrap(err, "Error while hashing"), 500
	}
	newUser := User{
		ID:       uint64(len(*db.Users)),
		Email:    (*input).Email,
		Username: (*input).Username,
		Password: string(passwordHashBytes),
	}
	*db.Users = append(*db.Users, newUser)
	return newUser, nil, 0
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func IsPasswordValid(password string) bool {
	return len(password) >= 6 && len(password) < 50
}

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	//TODO добавить аватар и в user
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
