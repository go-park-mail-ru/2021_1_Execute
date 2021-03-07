package api

import "github.com/labstack/echo"

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

type HttpError struct {
	Code uint
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
