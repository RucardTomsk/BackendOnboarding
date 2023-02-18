package model

import "github.com/RucardTomsk/BackendOnboarding/internal/domain/base"

type GetRoles struct {
	base.ResponseOK
	Roles []string
}
