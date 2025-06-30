package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"recommendation-service/internal/model"
)

type BookService struct {
	baseURL    string
	httpClient *http.Client
}

func NewBookService(baseURL string) *BookService {
	return &BookService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *BookService) GetAllBooks(pageSize int) ([]model.Book, error) {
	url := fmt.Sprintf("%s/api/books?page_size=%d", s.baseURL, pageSize)
	return s.getBooks(url)
}

func (s *BookService) GetBooksByCategory(category string, pageSize int) ([]model.Book, error) {
	url := fmt.Sprintf("%s/api/books/category/%s?page_size=%d", s.baseURL, category, pageSize)
	return s.getBooks(url)
}

func (s *BookService) GetBooksByAuthor(author string, pageSize int) ([]model.Book, error) {
	url := fmt.Sprintf("%s/api/books/author/%s?page_size=%d", s.baseURL, author, pageSize)
	return s.getBooks(url)
}

func (s *BookService) getBooks(url string) ([]model.Book, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("book service returned status %d", resp.StatusCode)
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
		if books, exists := dataMap["books"]; exists {
			booksJSON, _ := json.Marshal(books)
			var bookList []model.Book
			if err := json.Unmarshal(booksJSON, &bookList); err != nil {
				return nil, fmt.Errorf("failed to unmarshal books: %w", err)
			}
			return bookList, nil
		}
	}

	// Handle direct list response
	if bookList, ok := response.Data.([]interface{}); ok {
		booksJSON, _ := json.Marshal(bookList)
		var books []model.Book
		if err := json.Unmarshal(booksJSON, &books); err != nil {
			return nil, fmt.Errorf("failed to unmarshal books: %w", err)
		}
		return books, nil
	}

	return nil, fmt.Errorf("unexpected response format")
} 