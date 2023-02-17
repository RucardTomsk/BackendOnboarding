package model

type CreateUserRequest struct {
	UserName string
	Password string
	Email    string
}

type LoginUserRequest struct {
	Email    string
	Password string
}
