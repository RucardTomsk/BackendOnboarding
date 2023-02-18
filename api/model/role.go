package model

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
)

type GetRoles struct {
	base.ResponseOK
	Roles []string
}

type AddUserAndRoleRequest struct {
	UserID     string
	DivisionID string
	Role       string
}
