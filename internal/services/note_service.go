package services

import (
	"demo/internal/models"
	"demo/internal/repository"
)

type NoteService struct {
	noteRepo *repository.NoteRepository
	filmRepo *repository.FilmRepository
}

func NewNoteService(noteRepo *repository.NoteRepository, filmRepo *repository.FilmRepository) *NoteService {
	return &NoteService{noteRepo: noteRepo, filmRepo: filmRepo}
}

func (s *NoteService) Rate(userID, filmID uint, note int, commentaire string) (*models.Note, error) {
	if _, err := s.filmRepo.FindByID(filmID); err != nil {
		return nil, err
	}

	n := &models.Note{
		UserID:      userID,
		FilmID:      filmID,
		Note:        note,
		Commentaire: commentaire,
	}
	if err := s.noteRepo.Upsert(n); err != nil {
		return nil, err
	}

	return s.noteRepo.FindByUserAndFilm(userID, filmID)
}

func (s *NoteService) Get(userID, filmID uint) (*models.Note, error) {
	return s.noteRepo.FindByUserAndFilm(userID, filmID)
}

func (s *NoteService) List(userID uint) ([]models.Note, error) {
	return s.noteRepo.FindByUser(userID)
}

func (s *NoteService) Delete(userID, filmID uint) error {
	return s.noteRepo.Delete(userID, filmID)
}
