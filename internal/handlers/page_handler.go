package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"demo/internal/repository"
	"demo/internal/services"
	"demo/internal/views"
)

type PageHandler struct {
	filmService *services.FilmService
}

func NewPageHandler(filmService *services.FilmService) *PageHandler {
	return &PageHandler{filmService: filmService}
}

func (h *PageHandler) Films(c *gin.Context) {
	films, _, err := h.filmService.List(1, 50)
	if err != nil {
		c.String(http.StatusInternalServerError, "erreur interne")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = views.FilmsPage(films).Render(c.Request.Context(), c.Writer)
}

func (h *PageHandler) FilmDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "id invalide")
		return
	}

	film, err := h.filmService.Get(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Status(http.StatusNotFound)
			_ = views.FilmNotFoundPage().Render(c.Request.Context(), c.Writer)
			return
		}
		c.String(http.StatusInternalServerError, "erreur interne")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = views.FilmDetailPage(film).Render(c.Request.Context(), c.Writer)
}
