package db

import (
	"fmt"

	"agnos-test/internal/models"

	"gorm.io/gorm"
)

func Migrate(database *gorm.DB) error {
	if err := database.AutoMigrate(
		&models.Hospital{},
		&models.Staff{},
		&models.Patient{},
	); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
