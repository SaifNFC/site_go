package services

import (
	"time"

	"demo/internal/models"
	"demo/internal/repository"
)

type WatchlistService struct {
	watchlistRepo *repository.WatchlistRepository
	filmRepo      *repository.FilmRepository
}

func NewWatchlistService(watchlistRepo *repository.WatchlistRepository, filmRepo *repository.FilmRepository) *WatchlistService {
	return &WatchlistService{watchlistRepo: watchlistRepo, filmRepo: filmRepo}
}

func (s *WatchlistService) Add(userID, filmID uint) error {
	if _, err := s.filmRepo.FindByID(filmID); err != nil {
		return err
	}

	entry := &models.Watchlist{
		UserID:  userID,
		FilmID:  filmID,
		AddedAt: time.Now(),
	}
	return s.watchlistRepo.Add(entry)
}

func (s *WatchlistService) List(userID uint) ([]models.Watchlist, error) {
	return s.watchlistRepo.FindByUser(userID)
}

func (s *WatchlistService) Remove(userID, filmID uint) error {
	return s.watchlistRepo.Remove(userID, filmID)
}
