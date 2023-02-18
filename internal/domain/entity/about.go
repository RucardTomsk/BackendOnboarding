package entity

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/google/uuid"
)

type About struct {
	base.EntityWithGuidKey
	FIO         string `json:"fio"`
	Description string `json:"description"`
	Contact     string `json:"contact"`

	UserID uuid.UUID `json:"userID"`
}
