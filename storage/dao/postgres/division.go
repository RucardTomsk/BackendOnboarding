package postgresStorage

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DivisionStorage struct {
	db *gorm.DB
}

func NewDivisionStorage(db *gorm.DB) *DivisionStorage {
	return &DivisionStorage{db: db}
}

func (s DivisionStorage) Create(division *entity.Division, ctx context.Context) error {
	return s.db.Create(division).Error
}

func (s DivisionStorage) Retrieve(divisionID uuid.UUID, ctx context.Context) (*entity.Division, error) {
	var division entity.Division
	err := s.db.First(&division, divisionID).Error
	return &division, err
}

func (s DivisionStorage) Update(division *entity.Division, ctx context.Context) error {
	tx := s.db.Updates(division)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s DivisionStorage) Delete(divisionID uuid.UUID, ctx context.Context) error {
	tx := s.db.Delete(&entity.Division{}, divisionID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s DivisionStorage) Get() ([]entity.Division, error) {
	var division []entity.Division
	tx := s.db.Model(entity.Division{}).Find(&division)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return division, nil
}
