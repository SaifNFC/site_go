package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"demo/internal/config"
	"demo/internal/database"
	"demo/internal/repository"
	"demo/internal/services"
	"demo/internal/tmdb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database ping: %v", err)
	}
	log.Println("connexion à la base de données OK")

	filmRepo := repository.NewFilmRepository(db)
	tmdbClient := tmdb.NewClient(cfg.TMDBAPIKey)
	syncService := services.NewSyncService(tmdbClient, filmRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/sync/:tmdb_id", func(c *gin.Context) {
		tmdbID, err := strconv.Atoi(c.Param("tmdb_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tmdb_id invalide"})
			return
		}

		film, err := syncService.SyncFilm(tmdbID)
		if err != nil {
			if errors.Is(err, tmdb.ErrMovieNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Printf("sync film %d: %v", tmdbID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
			return
		}

		c.JSON(http.StatusOK, film)
	})

	router.Run(":" + cfg.Port)
}
