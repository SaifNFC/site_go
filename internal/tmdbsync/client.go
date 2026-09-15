package tmdbsync

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"demo/internal/models"
)

var ErrFilmNotFound = errors.New("film introuvable sur tmdb")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Sync(tmdbID int) (*models.Film, error) {
	endpoint := fmt.Sprintf("%s/sync/%d", c.baseURL, tmdbID)

	resp, err := c.httpClient.Post(endpoint, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("appel tmdb-sync: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrFilmNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb-sync a répondu avec le statut %d", resp.StatusCode)
	}

	var film models.Film
	if err := json.NewDecoder(resp.Body).Decode(&film); err != nil {
		return nil, fmt.Errorf("décodage réponse tmdb-sync: %w", err)
	}

	return &film, nil
}
