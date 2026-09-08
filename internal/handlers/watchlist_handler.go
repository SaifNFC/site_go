package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"demo/internal/middleware"
	"demo/internal/repository"
	"demo/internal/services"
)

type WatchlistHandler struct {
	watchlistService *services.WatchlistService
}

func NewWatchlistHandler(watchlistService *services.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{watchlistService: watchlistService}
}

func (h *WatchlistHandler) Add(c *gin.Context) {
	filmID, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint(middleware.UserIDKey)

	if err := h.watchlistService.Add(userID, filmID); err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *WatchlistHandler) List(c *gin.Context) {
	userID := c.GetUint(middleware.UserIDKey)

	entries, err := h.watchlistService.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"watchlist": entries})
}

func (h *WatchlistHandler) Remove(c *gin.Context) {
	filmID, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint(middleware.UserIDKey)

	if err := h.watchlistService.Remove(userID, filmID); err != nil {
		if errors.Is(err, repository.ErrWatchlistEntryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.Status(http.StatusNoContent)
}
