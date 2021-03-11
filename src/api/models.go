package api

import (
	"github.com/labstack/echo"
)

type Sessions map[string]int

type Database struct {
	echo.Context
	Users    *[]User
	Sessions *Sessions
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"-"`
}

type RegistrationResponse struct {
	ID int `json:"id"`
}

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserByIdResponse struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
}

type GetUserByIdBody struct {
	Response GetUserByIdResponse `json:"user"`
}

type PatchUserRequest struct {
	NewEmail    string `json:"email,omitempty"`
	NewUsername string `json:"username,omitempty"`
	NewPassword string `json:"password,omitempty"`
}
