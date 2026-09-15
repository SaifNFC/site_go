package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"demo/internal/tmdbsync"
)

type SyncHandler struct {
	syncClient *tmdbsync.Client
}

func NewSyncHandler(syncClient *tmdbsync.Client) *SyncHandler {
	return &SyncHandler{syncClient: syncClient}
}

func (h *SyncHandler) SyncFilm(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdb_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tmdb_id invalide"})
		return
	}

	film, err := h.syncClient.Sync(tmdbID)
	if err != nil {
		if errors.Is(err, tmdbsync.ErrFilmNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("sync film %d: %v", tmdbID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, film)
}
