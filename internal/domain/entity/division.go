package entity

import "github.com/RucardTomsk/BackendOnboarding/internal/domain/base"

type Division struct {
	base.EntityWithGuidKey
	Name        string
	Description string
}
