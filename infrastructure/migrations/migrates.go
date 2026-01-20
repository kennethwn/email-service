package migrations

import (
	"fmt"
	"worker-service/internal/dto"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MigrateAll(db *gorm.DB) {
	logrus.Info("Starting migrations...")
	err := db.AutoMigrate(
		&dto.EmailHistory{},
		&dto.ApiKey{},
	)
	if err != nil {
		logrus.Panic(fmt.Sprintf("failed to migrate all table, err: %v", err))
	}
	logrus.Info("Migration finished!")
}
