package tmdb

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestClient_GetMovie(t *testing.T) {
	_ = godotenv.Load("../../.env")

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY non défini, test ignoré")
	}

	client := NewClient(apiKey)

	movie, err := client.GetMovie(550) // Fight Club
	if err != nil {
		t.Fatalf("GetMovie a échoué: %v", err)
	}

	if movie.ID != 550 {
		t.Errorf("id attendu 550, obtenu %d", movie.ID)
	}
	if movie.Title == "" {
		t.Error("titre vide")
	}
}

func TestClient_GetMovie_NotFound(t *testing.T) {
	_ = godotenv.Load("../../.env")

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY non défini, test ignoré")
	}

	client := NewClient(apiKey)

	_, err := client.GetMovie(0)
	if err != ErrMovieNotFound {
		t.Errorf("erreur attendue ErrMovieNotFound, obtenu %v", err)
	}
}
