package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"demo/internal/models"
	"demo/internal/repository"
	"demo/internal/services"
)

type FilmHandler struct {
	filmService *services.FilmService
}

func NewFilmHandler(filmService *services.FilmService) *FilmHandler {
	return &FilmHandler{filmService: filmService}
}

type filmRequest struct {
	TMDBID    int    `json:"tmdb_id" binding:"required"`
	Titre     string `json:"titre" binding:"required"`
	Annee     int    `json:"annee"`
	PosterURL string `json:"poster_url"`
	Synopsis  string `json:"synopsis"`
}

func filmIDFromParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, errors.New("id invalide")
	}
	return uint(id), nil
}

func (h *FilmHandler) Create(c *gin.Context) {
	var req filmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	film := &models.Film{
		TMDBID:    req.TMDBID,
		Titre:     req.Titre,
		Annee:     req.Annee,
		PosterURL: req.PosterURL,
		Synopsis:  req.Synopsis,
	}

	if err := h.filmService.Create(film); err != nil {
		if errors.Is(err, services.ErrTMDBIDAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusCreated, film)
}

func (h *FilmHandler) Get(c *gin.Context) {
	id, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	film, err := h.filmService.Get(id)
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, film)
}

func (h *FilmHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	films, total, err := h.filmService.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"films": films, "total": total, "page": page})
}

func (h *FilmHandler) Update(c *gin.Context) {
	id, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req filmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	film := &models.Film{
		ID:        id,
		TMDBID:    req.TMDBID,
		Titre:     req.Titre,
		Annee:     req.Annee,
		PosterURL: req.PosterURL,
		Synopsis:  req.Synopsis,
	}

	if err := h.filmService.Update(film); err != nil {
		switch {
		case errors.Is(err, repository.ErrFilmNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTMDBIDAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		}
		return
	}

	c.JSON(http.StatusOK, film)
}

func (h *FilmHandler) Delete(c *gin.Context) {
	id, err := filmIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.filmService.Delete(id); err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.Status(http.StatusNoContent)
}
