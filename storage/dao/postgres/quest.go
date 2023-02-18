package postgresStorage

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestStorage struct {
	db *gorm.DB
}

func NewQuestStorage(db *gorm.DB) *QuestStorage {
	return &QuestStorage{
		db: db,
	}
}

func (s QuestStorage) CreateQuest(quest *entity.Quest, ctx context.Context) error {
	return s.db.Create(quest).Error
}

func (s QuestStorage) CreateStage(stage *entity.Stage, ctx context.Context) error {
	return s.db.Create(stage).Error
}

func (s QuestStorage) RetrieveQuest(questID uuid.UUID, ctx context.Context) (*entity.Quest, error) {
	var quest entity.Quest
	err := s.db.Preload("Stages").First(&quest, questID).Error
	return &quest, err
}

func (s QuestStorage) UpdateQuest(quest *entity.Quest, ctx context.Context) error {
	tx := s.db.Updates(quest)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s QuestStorage) UpdateStage(stage *entity.Quest, ctx context.Context) error {
	tx := s.db.Updates(stage)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s QuestStorage) GetQuest() ([]entity.Quest, error) {
	var quests []entity.Quest
	tx := s.db.Preload("Stages").Model(entity.Quest{}).Find(&quests)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return quests, nil
}
