package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"demo/internal/models"
)

var ErrNoteNotFound = errors.New("note introuvable")

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Upsert crée la note si l'utilisateur n'a pas encore noté ce film, ou la remplace sinon.
// S'appuie sur ON CONFLICT (PostgreSQL) plutôt qu'un SELECT puis INSERT/UPDATE séparés.
func (r *NoteRepository) Upsert(note *models.Note) error {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "film_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"note", "commentaire"}),
	}).Create(note).Error
	if err != nil {
		return fmt.Errorf("upsert note: %w", err)
	}
	return nil
}

func (r *NoteRepository) FindByUserAndFilm(userID, filmID uint) (*models.Note, error) {
	var note models.Note
	if err := r.db.Where("user_id = ? AND film_id = ?", userID, filmID).First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, fmt.Errorf("recherche note: %w", err)
	}
	return &note, nil
}

func (r *NoteRepository) FindByUser(userID uint) ([]models.Note, error) {
	var notes []models.Note
	if err := r.db.Preload("Film").Where("user_id = ?", userID).Order("created_at desc").Find(&notes).Error; err != nil {
		return nil, fmt.Errorf("liste notes: %w", err)
	}
	return notes, nil
}

func (r *NoteRepository) Delete(userID, filmID uint) error {
	result := r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&models.Note{})
	if result.Error != nil {
		return fmt.Errorf("suppression note: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
