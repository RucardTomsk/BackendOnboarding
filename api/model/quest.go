package model

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/google/uuid"
)

type CreateQuestRequest struct {
	DivisionID  string
	Name        string
	Description string
}

type CreateStageRequest struct {
	QuestID     string
	Name        string
	Description string
}

type StageObject struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type QuestObject struct {
	ID          uuid.UUID
	Name        string
	Description string
	Stages      []StageObject
}

type GetAllQuestResponse struct {
	base.ResponseOK
	Quests []QuestObject
}
