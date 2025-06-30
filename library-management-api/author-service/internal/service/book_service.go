package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"author-service/internal/model"
)

// BookService book microservice ile iletişim interface'i
type BookService interface {
	GetBooksByAuthor(authorName string) ([]model.BookInfo, error)
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

// GetBooksByAuthor book service'den yazar kitaplarını getirir
func (s *HTTPBookService) GetBooksByAuthor(authorName string) ([]model.BookInfo, error) {
	url := fmt.Sprintf("%s/api/books/author/%s", s.baseURL, authorName)
	
	log.Printf("Book service'e istek gönderiliyor: %s", url)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, model.NewAuthorError("BOOK_SERVICE_ERROR", "Book service'e bağlanamadı", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, model.NewAuthorError("BOOK_SERVICE_HTTP_ERROR", fmt.Sprintf("Book service'den hata: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewAuthorError("BOOK_SERVICE_READ_ERROR", "Book service yanıtı okunamadı", err)
	}

	var bookResponse struct {
		Data []model.BookInfo `json:"data"`
	}
	
	if err := json.Unmarshal(body, &bookResponse); err != nil {
		return nil, model.NewAuthorError("BOOK_SERVICE_PARSE_ERROR", "Book service yanıtı parse edilemedi", err)
	}

	return bookResponse.Data, nil
} 