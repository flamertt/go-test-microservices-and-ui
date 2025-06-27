package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"recommendation-service/internal/model"
)

type GenreService struct {
	baseURL    string
	httpClient *http.Client
}

func NewGenreService(baseURL string) *GenreService {
	return &GenreService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *GenreService) GetAllGenres(pageSize int) ([]model.Genre, error) {
	url := fmt.Sprintf("%s/genres?page_size=%d", s.baseURL, pageSize)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch genres: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("genre service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Data interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Handle paginated response
	if dataMap, ok := response.Data.(map[string]interface{}); ok {
		if genres, exists := dataMap["genres"]; exists {
			genresJSON, _ := json.Marshal(genres)
			var genreList []model.Genre
			if err := json.Unmarshal(genresJSON, &genreList); err != nil {
				return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
			}
			return genreList, nil
		}
	}

	return nil, fmt.Errorf("unexpected response format")
} 