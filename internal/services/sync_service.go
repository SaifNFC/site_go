package services

import (
	"errors"
	"strconv"

	"demo/internal/models"
	"demo/internal/repository"
	"demo/internal/tmdb"
)

type SyncService struct {
	tmdbClient *tmdb.Client
	filmRepo   *repository.FilmRepository
}

func NewSyncService(tmdbClient *tmdb.Client, filmRepo *repository.FilmRepository) *SyncService {
	return &SyncService{tmdbClient: tmdbClient, filmRepo: filmRepo}
}

func (s *SyncService) SyncFilm(tmdbID int) (*models.Film, error) {
	movie, err := s.tmdbClient.GetMovie(tmdbID)
	if err != nil {
		return nil, err
	}

	annee := 0
	if len(movie.ReleaseDate) >= 4 {
		annee, _ = strconv.Atoi(movie.ReleaseDate[:4])
	}

	existing, err := s.filmRepo.FindByTMDBID(tmdbID)
	if err != nil && !errors.Is(err, repository.ErrFilmNotFound) {
		return nil, err
	}

	if existing == nil {
		film := &models.Film{
			TMDBID:    tmdbID,
			Titre:     movie.Title,
			Annee:     annee,
			PosterURL: movie.PosterPath,
			Synopsis:  movie.Overview,
		}
		if err := s.filmRepo.Create(film); err != nil {
			return nil, err
		}
		return film, nil
	}

	existing.Titre = movie.Title
	existing.Annee = annee
	existing.PosterURL = movie.PosterPath
	existing.Synopsis = movie.Overview
	if err := s.filmRepo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
