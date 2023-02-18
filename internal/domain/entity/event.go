package entity

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"time"
)

type Event struct {
	base.EntityWithGuidKey
	DataTime    time.Time `json:"dataTime"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Format      string    `json:"format"`
}
