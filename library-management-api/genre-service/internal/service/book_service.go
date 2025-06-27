package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"genre-service/internal/model"
)

// BookService book microservice ile iletişim interface'i
type BookService interface {
	GetBooksByCategory(categoryName string) ([]model.BookInfo, error)
	GetBooksByCategoryWithPagination(categoryName string, page, pageSize int) ([]model.BookInfo, error)
	GetBookCountByCategory(categoryName string) (int, error)
}

// HTTPBookService HTTP üzerinden book service implementasyonu
type HTTPBookService struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPBookService yeni HTTP book service oluşturur
func NewHTTPBookService(baseURL string) BookService {
	return &HTTPBookService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBooksByCategory book service'den kategori kitaplarını getirir
func (s *HTTPBookService) GetBooksByCategory(categoryName string) ([]model.BookInfo, error) {
	url := fmt.Sprintf("%s/books/category/%s", s.baseURL, categoryName)
	
	log.Printf("Book service'e istek gönderiliyor: %s", url)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, model.NewGenreError("BOOK_SERVICE_ERROR", "Book service'e bağlanamadı", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, model.NewGenreError("BOOK_SERVICE_HTTP_ERROR", fmt.Sprintf("Book service'den hata: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewGenreError("BOOK_SERVICE_READ_ERROR", "Book service yanıtı okunamadı", err)
	}

	var bookResponse struct {
		Data []model.BookInfo `json:"data"`
	}
	
	if err := json.Unmarshal(body, &bookResponse); err != nil {
		return nil, model.NewGenreError("BOOK_SERVICE_PARSE_ERROR", "Book service yanıtı parse edilemedi", err)
	}

	return bookResponse.Data, nil
}

// GetBooksByCategoryWithPagination book service'den kategori kitaplarını sayfalanmış olarak getirir
func (s *HTTPBookService) GetBooksByCategoryWithPagination(categoryName string, page, pageSize int) ([]model.BookInfo, error) {
	url := fmt.Sprintf("%s/books/category/%s?page=%d&page_size=%d", s.baseURL, categoryName, page, pageSize)
	
	log.Printf("Book service'e sayfalanmış istek gönderiliyor: %s", url)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, model.NewGenreError("BOOK_SERVICE_ERROR", "Book service'e bağlanamadı", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, model.NewGenreError("BOOK_SERVICE_HTTP_ERROR", fmt.Sprintf("Book service'den hata: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewGenreError("BOOK_SERVICE_READ_ERROR", "Book service yanıtı okunamadı", err)
	}

	var response struct {
		Data struct {
			Books []model.BookInfo `json:"books"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, model.NewGenreError("BOOK_SERVICE_PARSE_ERROR", "Book service yanıtı parse edilemedi", err)
	}

	return response.Data.Books, nil
}

// GetBookCountByCategory book service'den kategori kitap sayısını getirir
func (s *HTTPBookService) GetBookCountByCategory(categoryName string) (int, error) {
	// Pagination ile 1 sayfa, 1 eleman isteyerek total count'u al (optimize edilmiş)
	url := fmt.Sprintf("%s/books/category/%s?page=1&page_size=1", s.baseURL, categoryName)
	
	log.Printf("Book service'e count isteği gönderiliyor: %s", url)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return 0, model.NewGenreError("BOOK_SERVICE_ERROR", "Book service'e bağlanamadı", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, model.NewGenreError("BOOK_SERVICE_HTTP_ERROR", fmt.Sprintf("Book service'den hata: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, model.NewGenreError("BOOK_SERVICE_READ_ERROR", "Book service yanıtı okunamadı", err)
	}

	var response struct {
		Data struct {
			Total int `json:"total"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, model.NewGenreError("BOOK_SERVICE_PARSE_ERROR", "Book service yanıtı parse edilemedi", err)
	}

	return response.Data.Total, nil
} 