package base

import (
	"github.com/google/uuid"
	"time"
)

// EntityWithGuidKey is a base DB entity with uuid.UUID as a primary key.
type EntityWithGuidKey struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v1();primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

// ArchivableEntityWithGuidKey is an EntityWithGuidKey struct with extra ArchivedAt field.
type ArchivableEntityWithGuidKey struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v1();primaryKey"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	DeletedAt  *time.Time `json:"-"`
}

// EntityWithIntegerKey is a base DB entity with uint as a primary key.
type EntityWithIntegerKey struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}
