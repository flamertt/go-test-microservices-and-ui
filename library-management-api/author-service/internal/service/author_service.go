package service

import (
	"log"

	"author-service/internal/model"
	"author-service/internal/repository"
)

// AuthorService yazar iş mantığı interface'i
type AuthorService interface {
	GetPaginatedAuthors(params *model.AuthorSearchParams) (*model.PaginatedAuthors, error)
	GetAuthorByName(name string) ([]model.Author, error)
	GetEnrichedAuthorByName(name string) (*model.EnrichedAuthor, error)
}

// AuthorServiceImpl AuthorService implementasyonu
type AuthorServiceImpl struct {
	authorRepo  repository.AuthorRepository
	bookService BookService
}

// NewAuthorService yeni author service oluşturur
func NewAuthorService(authorRepo repository.AuthorRepository, bookService BookService) AuthorService {
	return &AuthorServiceImpl{
		authorRepo:  authorRepo,
		bookService: bookService,
	}
}

// GetPaginatedAuthors sayfalı yazar listesi getirir
func (s *AuthorServiceImpl) GetPaginatedAuthors(params *model.AuthorSearchParams) (*model.PaginatedAuthors, error) {
	// Parametreleri doğrula
	if err := s.validateSearchParams(params); err != nil {
		return nil, err
	}

	return s.authorRepo.GetPaginatedAuthors(params)
}

// GetAuthorByName isim ile yazar arama
func (s *AuthorServiceImpl) GetAuthorByName(name string) ([]model.Author, error) {
	if name == "" {
		return nil, model.ErrInvalidAuthorName
	}

	return s.authorRepo.GetAuthorByName(name)
}

// GetEnrichedAuthorByName kitap bilgisiyle zenginleştirilmiş yazar getirir
func (s *AuthorServiceImpl) GetEnrichedAuthorByName(name string) (*model.EnrichedAuthor, error) {
	// Önce yazar bilgisini al
	authors, err := s.GetAuthorByName(name)
	if err != nil {
		return nil, err
	}

	if len(authors) == 0 {
		return nil, model.ErrAuthorNotFound
	}

	// İlk yazarı al (genelde tek olur)
	author := authors[0]

	// Book service'den bu yazarın kitaplarını al
	books, err := s.bookService.GetBooksByAuthor(name)
	if err != nil {
		log.Printf("Yazar kitapları alınamadı: %v", err)
		// Hata olsa bile yazar bilgisini döndür
		books = []model.BookInfo{}
	}

	return author.ToEnriched(books), nil
}

// validateSearchParams arama parametrelerini doğrular
func (s *AuthorServiceImpl) validateSearchParams(params *model.AuthorSearchParams) error {
	if params.Page < 1 {
		return model.ErrInvalidPage
	}
	
	if params.PageSize < 1 || params.PageSize > 100 {
		return model.ErrInvalidPageSize
	}

	return nil
} 