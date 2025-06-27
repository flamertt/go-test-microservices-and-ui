package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"book-service/internal/model"
)

// AuthorService author microservice ile iletişim interface'i
type AuthorService interface {
	GetAuthorInfo(authorName string) (*model.AuthorInfo, error)
}

// HTTPAuthorService HTTP üzerinden author service implementasyonu
type HTTPAuthorService struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPAuthorService yeni HTTP author service oluşturur
func NewHTTPAuthorService(baseURL string) AuthorService {
	return &HTTPAuthorService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAuthorInfo author service'den yazar bilgisini getirir
func (s *HTTPAuthorService) GetAuthorInfo(authorName string) (*model.AuthorInfo, error) {
	url := fmt.Sprintf("%s/authors/search?name=%s", s.baseURL, authorName)
	
	log.Printf("Author service'e istek gönderiliyor: %s", url)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, model.NewBookError("AUTHOR_SERVICE_ERROR", "Author service'e bağlanamadı", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, model.NewBookError("AUTHOR_SERVICE_HTTP_ERROR", fmt.Sprintf("Author service'den hata: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewBookError("AUTHOR_SERVICE_READ_ERROR", "Author service yanıtı okunamadı", err)
	}

	var authorResponse struct {
		Data []map[string]interface{} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &authorResponse); err != nil {
		return nil, model.NewBookError("AUTHOR_SERVICE_PARSE_ERROR", "Author service yanıtı parse edilemedi", err)
	}

	if len(authorResponse.Data) == 0 {
		return &model.AuthorInfo{
			Name:      authorName,
			Biography: "Yazar bilgisi bulunamadı",
		}, nil
	}

	// İlk yazarı al
	author := authorResponse.Data[0]
	
	return &model.AuthorInfo{
		Name:      authorName,
		Biography: fmt.Sprintf("Yazar hakkında bilgi: %v", author),
	}, nil
} 