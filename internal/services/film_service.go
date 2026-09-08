package services

import (
	"errors"

	"demo/internal/models"
	"demo/internal/repository"
)

var ErrTMDBIDAlreadyExists = errors.New("un film avec cet identifiant TMDB existe déjà")

const (
	defaultLimit = 20
	maxLimit     = 100
)

type FilmService struct {
	filmRepo *repository.FilmRepository
}

func NewFilmService(filmRepo *repository.FilmRepository) *FilmService {
	return &FilmService{filmRepo: filmRepo}
}

func (s *FilmService) Create(film *models.Film) error {
	if _, err := s.filmRepo.FindByTMDBID(film.TMDBID); err == nil {
		return ErrTMDBIDAlreadyExists
	} else if !errors.Is(err, repository.ErrFilmNotFound) {
		return err
	}

	return s.filmRepo.Create(film)
}

func (s *FilmService) Get(id uint) (*models.Film, error) {
	return s.filmRepo.FindByID(id)
}

func (s *FilmService) List(page, limit int) ([]models.Film, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > maxLimit {
		limit = defaultLimit
	}
	offset := (page - 1) * limit

	return s.filmRepo.FindAll(limit, offset)
}

func (s *FilmService) Update(film *models.Film) error {
	if _, err := s.filmRepo.FindByID(film.ID); err != nil {
		return err
	}

	existing, err := s.filmRepo.FindByTMDBID(film.TMDBID)
	if err == nil && existing.ID != film.ID {
		return ErrTMDBIDAlreadyExists
	} else if err != nil && !errors.Is(err, repository.ErrFilmNotFound) {
		return err
	}

	return s.filmRepo.Update(film)
}

func (s *FilmService) Delete(id uint) error {
	return s.filmRepo.Delete(id)
}
