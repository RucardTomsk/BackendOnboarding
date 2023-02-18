package migration

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Division{},
		&entity.Quest{},
		&entity.Stage{},
		&entity.User{},
		&entity.About{},
		&entity.Event{},
	)
}
