package model

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/google/uuid"
)

type CreateDivisionRequest struct {
	Name        string
	Description string
}

type DivisionObject struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type GetDivisions struct {
	base.ResponseOK
	Divisions []DivisionObject
}
