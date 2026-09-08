package tmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.themoviedb.org/3"

var ErrMovieNotFound = errors.New("film introuvable sur tmdb")

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
	Overview    string `json:"overview"`
}

func (c *Client) GetMovie(tmdbID int) (*Movie, error) {
	endpoint := fmt.Sprintf("%s/movie/%d?api_key=%s&language=fr-FR", baseURL, tmdbID, url.QueryEscape(c.apiKey))

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("appel tmdb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrMovieNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb a répondu avec le statut %d", resp.StatusCode)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("décodage réponse tmdb: %w", err)
	}

	return &movie, nil
}
