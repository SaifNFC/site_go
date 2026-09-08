package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"demo/internal/config"
	"demo/internal/models"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connexion à la base de données: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Film{},
		&models.Note{},
		&models.Watchlist{},
	); err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	return nil
}
