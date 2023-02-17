package entity

import "github.com/RucardTomsk/BackendOnboarding/internal/domain/base"

type User struct {
	base.EntityWithGuidKey
	UserName   string `json:"userName" gorm:"uniqueIndex"`
	Password   string `json:"password"`
	Email      string `json:"email" gorm:"uniqueIndex"`
	TelegramID string `json:"tgID"`
}
