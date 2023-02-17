package postgresStorage

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s UserStorage) Create(user *entity.User, ctx context.Context) error {
	return s.db.Create(user).Error
}

func (s UserStorage) Retrieve(userID uuid.UUID, ctx context.Context) (*entity.User, error) {
	var user entity.User
	err := s.db.First(&user, userID).Error
	return &user, err
}

func (s UserStorage) RetrieveTo(Email string, Password string, ctx context.Context) (*entity.User, error) {
	var user entity.User
	err := s.db.Model(entity.User{}).
		Where("email = ?", Email).
		Where("password = ?", Password).
		First(&user).Error

	return &user, err
}

func (s UserStorage) Update(user *entity.User, ctx context.Context) error {
	tx := s.db.Updates(user)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s UserStorage) Delete(userID uuid.UUID, ctx context.Context) error {
	tx := s.db.Delete(&entity.User{}, userID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s UserStorage) Get() ([]entity.User, error) {
	var users []entity.User
	tx := s.db.Model(entity.User{}).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return users, nil
}
