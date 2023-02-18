package model

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
)

type CreateUserRequest struct {
	Password string
	Email    string
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type UserInfoResponse struct {
	base.ResponseOK
	Email string
}

type UserObject struct {
	Email string
	About AboutObject
}

type GetUsersRequest struct {
	base.ResponseOK
	Users []UserObject
}
