package entity

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/google/uuid"
)

type Quest struct {
	base.EntityWithGuidKey
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`

	Stages []Stage `json:"stages" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	DivisionID uuid.UUID
	Division   *Division `gorm:"-:all"`
}

type Stage struct {
	base.EntityWithGuidKey
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`

	QuestID uuid.UUID
	Quest   *Quest `gorm:"-:all"`
}
