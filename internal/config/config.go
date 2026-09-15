package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	TMDBAPIKey  string
	TMDBSyncURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmdbSyncURL := os.Getenv("TMDB_SYNC_URL")
	if tmdbSyncURL == "" {
		tmdbSyncURL = "http://localhost:8082"
	}

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		TMDBAPIKey:  os.Getenv("TMDB_API_KEY"),
		TMDBSyncURL: tmdbSyncURL,
	}, nil
}
