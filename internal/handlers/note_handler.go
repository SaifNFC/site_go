package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"demo/internal/middleware"
	"demo/internal/repository"
	"demo/internal/services"
)

type NoteHandler struct {
	noteService *services.NoteService
}

func NewNoteHandler(noteService *services.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

type noteRequest struct {
	Note        int    `json:"note" binding:"required,min=1,max=10"`
	Commentaire string `json:"commentaire"`
}

func (h *NoteHandler) Rate(c *gin.Context) {
	filmID, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint(middleware.UserIDKey)

	note, err := h.noteService.Rate(userID, filmID, req.Note, req.Commentaire)
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Get(c *gin.Context) {
	filmID, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint(middleware.UserIDKey)

	note, err := h.noteService.Get(userID, filmID)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) List(c *gin.Context) {
	userID := c.GetUint(middleware.UserIDKey)

	notes, err := h.noteService.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (h *NoteHandler) Delete(c *gin.Context) {
	filmID, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint(middleware.UserIDKey)

	if err := h.noteService.Delete(userID, filmID); err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.Status(http.StatusNoContent)
}
