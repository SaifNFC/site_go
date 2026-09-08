package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"demo/internal/models"
)

var ErrFilmNotFound = errors.New("film introuvable")

type FilmRepository struct {
	db *gorm.DB
}

func NewFilmRepository(db *gorm.DB) *FilmRepository {
	return &FilmRepository{db: db}
}

func (r *FilmRepository) Create(film *models.Film) error {
	if err := r.db.Create(film).Error; err != nil {
		return fmt.Errorf("création film: %w", err)
	}
	return nil
}

func (r *FilmRepository) FindByID(id uint) (*models.Film, error) {
	var film models.Film
	if err := r.db.First(&film, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFilmNotFound
		}
		return nil, fmt.Errorf("recherche film par id: %w", err)
	}
	return &film, nil
}

func (r *FilmRepository) FindByTMDBID(tmdbID int) (*models.Film, error) {
	var film models.Film
	if err := r.db.Where("tmdb_id = ?", tmdbID).First(&film).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFilmNotFound
		}
		return nil, fmt.Errorf("recherche film par tmdb_id: %w", err)
	}
	return &film, nil
}

func (r *FilmRepository) FindAll(limit, offset int) ([]models.Film, int64, error) {
	var films []models.Film
	var total int64

	if err := r.db.Model(&models.Film{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("comptage films: %w", err)
	}

	if err := r.db.Limit(limit).Offset(offset).Order("id").Find(&films).Error; err != nil {
		return nil, 0, fmt.Errorf("liste films: %w", err)
	}

	return films, total, nil
}

func (r *FilmRepository) Update(film *models.Film) error {
	if err := r.db.Save(film).Error; err != nil {
		return fmt.Errorf("mise à jour film: %w", err)
	}
	return nil
}

func (r *FilmRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Film{}, id)
	if result.Error != nil {
		return fmt.Errorf("suppression film: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrFilmNotFound
	}
	return nil
}
