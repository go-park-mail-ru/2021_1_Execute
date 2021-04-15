package models

import "2021_1_Execute/internal/users"

func LoginRequestToUser(input *UserLoginRequest) users.User {
	return users.User{
		Email:    input.Email,
		Password: input.Password,
	}
}

func RegistrationRequestToUser(input *UserRegistrationRequest) users.User {
	return users.User{
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	}
}

type UserLoginRequest struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"password" `
}

type EntranceResponse struct {
	ID int `json:"id"`
}

type UserRegistrationRequest struct {
	Email    string `json:"email" valid:"email"`
	Username string `json:"username" valid:"username"`
	Password string `json:"password" valid:"password" `
}

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

func CreateGetUserByIdBody(user users.User) GetUserByIdResponse {
	return GetUserByIdResponse{
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.Avatar,
	}
}

func CreateGetUserByIdResponse(user users.User) GetUserByIdBody {
	return GetUserByIdBody{
		Response: CreateGetUserByIdBody(user),
	}
}

func CreateUserFromPatchRequest(input *PatchUserRequest) users.User {
	return users.User{
		Email:    input.NewEmail,
		Username: input.NewUsername,
		Password: input.NewPassword,
		Avatar:   "",
	}
}
