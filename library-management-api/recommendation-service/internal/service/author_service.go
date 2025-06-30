package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"recommendation-service/internal/model"
)

type AuthorService struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthorService(baseURL string) *AuthorService {
	return &AuthorService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *AuthorService) GetAllAuthors(pageSize int) ([]model.Author, error) {
	url := fmt.Sprintf("%s/api/authors?page_size=%d", s.baseURL, pageSize)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("author service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Data struct {
			Authors []model.Author `json:"authors"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Data.Authors, nil
} 