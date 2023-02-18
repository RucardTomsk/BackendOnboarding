package model

import "github.com/RucardTomsk/BackendOnboarding/internal/domain/base"

type CreateDivisionRequest struct {
	Name        string
	Description string
}

type DivisionObject struct {
	Name        string
	Description string
}

type GetDivisions struct {
	base.ResponseOK
	Divisions []DivisionObject
}
