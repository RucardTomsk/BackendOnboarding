package entity

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/google/uuid"
)

type User struct {
	base.EntityWithGuidKey
	Password   string `json:"password"`
	Email      string `json:"email" gorm:"uniqueIndex"`
	TelegramID string `json:"tgID"`

	Points int `json:"points"`

	AboutID uuid.UUID `json:"aboutID"`
	About   *About    `json:"about,omitempty" gorm:"-,all"`
}
