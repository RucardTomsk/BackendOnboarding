package postgresStorage

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"gorm.io/gorm"
)

type AboutStorage struct {
	db *gorm.DB
}

func NewAboutStorage(db *gorm.DB) *AboutStorage {
	return &AboutStorage{db: db}
}

func (s AboutStorage) Create(about *entity.About, ctx context.Context) error {
	return s.db.Create(about).Error
}

func (s AboutStorage) Update(about *entity.About, ctx context.Context) error {
	tx := s.db.Updates(about)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
