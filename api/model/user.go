package model

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/google/uuid"
)

const (
	AdminEmail    = "admin@example.com"
	AdminPassword = "admin"
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
	User UserObject
}

type UserObject struct {
	ID    uuid.UUID
	Email string
	About AboutObject
}

type OverQuestObject struct {
	ID          uuid.UUID
	Name        string
	Description string
	Stages      []StageObject
}
type UserQuestObject struct {
	Division DivisionObject
	Quests   []OverQuestObject
}

type AllQuestUserResponse struct {
	base.ResponseOK
	DivQuests []UserQuestObject
}

type GetUsersRequest struct {
	base.ResponseOK
	Users []UserObject
}
