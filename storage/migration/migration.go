package migration

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
