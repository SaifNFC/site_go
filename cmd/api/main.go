package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"demo/internal/config"
	"demo/internal/database"
	"demo/internal/handlers"
	"demo/internal/middleware"
	"demo/internal/repository"
	"demo/internal/services"
	"demo/internal/tmdbsync"
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

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database: %v", err)
	}
	log.Println("migrations OK")

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	filmRepo := repository.NewFilmRepository(db)
	filmService := services.NewFilmService(filmRepo)
	filmHandler := handlers.NewFilmHandler(filmService)

	noteRepo := repository.NewNoteRepository(db)
	noteService := services.NewNoteService(noteRepo, filmRepo)
	noteHandler := handlers.NewNoteHandler(noteService)

	watchlistRepo := repository.NewWatchlistRepository(db)
	watchlistService := services.NewWatchlistService(watchlistRepo, filmRepo)
	watchlistHandler := handlers.NewWatchlistHandler(watchlistService)

	pageHandler := handlers.NewPageHandler(filmService)

	syncClient := tmdbsync.NewClient(cfg.TMDBSyncURL)
	syncHandler := handlers.NewSyncHandler(syncClient)

	router := gin.Default()

	router.Static("/static", "./web/static")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/", pageHandler.Films)
	router.GET("/films/:id/view", pageHandler.FilmDetail)

	auth := router.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	protected := router.Group("/")
	protected.Use(middleware.Auth(authService))
	protected.GET("/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": c.GetUint(middleware.UserIDKey)})
	})

	router.GET("/films", filmHandler.List)
	router.GET("/films/:id", filmHandler.Get)

	films := protected.Group("/films")
	films.POST("", filmHandler.Create)
	films.PUT("/:id", filmHandler.Update)
	films.DELETE("/:id", filmHandler.Delete)

	protected.PUT("/films/:id/note", noteHandler.Rate)
	protected.GET("/films/:id/note", noteHandler.Get)
	protected.DELETE("/films/:id/note", noteHandler.Delete)
	protected.GET("/me/notes", noteHandler.List)

	protected.POST("/films/:id/watchlist", watchlistHandler.Add)
	protected.DELETE("/films/:id/watchlist", watchlistHandler.Remove)
	protected.GET("/me/watchlist", watchlistHandler.List)

	protected.POST("/admin/films/:tmdb_id/sync", syncHandler.SyncFilm)

	router.Run(":" + cfg.Port)
}
