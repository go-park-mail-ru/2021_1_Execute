package api

import (
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func (db *Database) CreateUser(input *UserRegistrationRequest) (User, error) {
	if !IsEmailValid(input.Email) {
		return User{}, BadRequestError
	}

	if !IsPasswordValid((*input).Password) {
		return User{}, BadRequestError
	}

	for _, user := range *db.Users {
		if user.Email == input.Email {
			return User{}, ConflictError
		}
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, InternalServerError
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

func (db *Database) IsEmailUniq(userID int, email string) bool {
	for _, user := range *db.Users {
		if userID != user.ID && user.Email == email {
			return false
		}
	}
	return true
}

func (db *Database) UpdateUser(userID int, username, email, password, avatar string) error {
	switch {
	case email != "" && !IsEmailValid(email):
		return BadRequestError
	case email != "" && !db.IsEmailUniq(userID, email):
		return ConflictError
	case password != "" && !IsPasswordValid(password):
		return BadRequestError
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
					return InternalServerError
				}
				(*db.Users)[i].Password = string(hash)
			}
			if avatar != "" {
				if user.Avatar != "" {
					err := deleteFile(user.Avatar)
					if err != nil {
						return InternalServerError
					}
				}
				(*db.Users)[i].Avatar = avatar
			}
			return nil
		}
	}

	return NotFoundError
}

func (db *Database) DeleteUser(userID int) error {
	for i, user := range *db.Users {
		if user.ID == userID {
			*db.Users = append((*db.Users)[:i], (*db.Users)[i+1:]...)
			return nil
		}
	}
	return NotFoundError
}
