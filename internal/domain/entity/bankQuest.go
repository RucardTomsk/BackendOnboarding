package entity

import "github.com/RucardTomsk/BackendOnboarding/internal/domain/base"

type BankQuest struct {
	base.ResponseOKWithGUID
	Name        string
	Description string
	IsDone      bool

	Stages []BankStage
}

type BankStage struct {
	base.ResponseOKWithGUID
	Name        string
	Description string
	IsDone      bool
}
