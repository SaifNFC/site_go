package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"demo/internal/models"
)

var ErrWatchlistEntryNotFound = errors.New("film absent de la watchlist")

type WatchlistRepository struct {
	db *gorm.DB
}

func NewWatchlistRepository(db *gorm.DB) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}

// Add ignore silencieusement si le film est déjà dans la watchlist (idempotent).
func (r *WatchlistRepository) Add(entry *models.Watchlist) error {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "film_id"}},
		DoNothing: true,
	}).Create(entry).Error
	if err != nil {
		return fmt.Errorf("ajout watchlist: %w", err)
	}
	return nil
}

func (r *WatchlistRepository) FindByUser(userID uint) ([]models.Watchlist, error) {
	var entries []models.Watchlist
	if err := r.db.Preload("Film").Where("user_id = ?", userID).Order("added_at desc").Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("liste watchlist: %w", err)
	}
	return entries, nil
}

func (r *WatchlistRepository) Remove(userID, filmID uint) error {
	result := r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&models.Watchlist{})
	if result.Error != nil {
		return fmt.Errorf("suppression watchlist: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrWatchlistEntryNotFound
	}
	return nil
}
