package service

import (
	"errors"
	"log"

	"book-service/internal/model"
	"book-service/internal/repository"
)

// BookService kitap iş mantığı interface'i
type BookService interface {
	GetPaginatedBooks(params *model.BookSearchParams) (*model.PaginatedBooks, error)
	GetBookByID(id int) (*model.Book, error)
	GetEnrichedBookByID(id int) (*model.EnrichedBook, error)
	GetBooksByAuthor(authorName string) ([]model.Book, error)
	GetBooksByCategory(categoryName string) ([]model.Book, error)
	GetBooksByCategoryWithPagination(categoryName string, params *model.BookSearchParams) (*model.PaginatedBooks, error)
	GetEnrichedBooks(params *model.BookSearchParams) ([]*model.EnrichedBook, error)
}

// BookServiceImpl BookService implementasyonu
type BookServiceImpl struct {
	bookRepo      repository.BookRepository
	authorService AuthorService
}

// NewBookService yeni book service oluşturur
func NewBookService(bookRepo repository.BookRepository, authorService AuthorService) BookService {
	return &BookServiceImpl{
		bookRepo:      bookRepo,
		authorService: authorService,
	}
}

// GetPaginatedBooks sayfalı kitap listesi getirir
func (s *BookServiceImpl) GetPaginatedBooks(params *model.BookSearchParams) (*model.PaginatedBooks, error) {
	// Parametreleri doğrula
	if err := s.validateSearchParams(params); err != nil {
		return nil, err
	}

	return s.bookRepo.GetPaginatedBooks(params)
}

// GetBookByID ID'ye göre kitap getirir
func (s *BookServiceImpl) GetBookByID(id int) (*model.Book, error) {
	if id <= 0 {
		return nil, model.ErrInvalidBookID
	}

	return s.bookRepo.GetBookByID(id)
}

// GetEnrichedBookByID ID'ye göre yazar bilgisiyle zenginleştirilmiş kitap getirir
func (s *BookServiceImpl) GetEnrichedBookByID(id int) (*model.EnrichedBook, error) {
	// Önce kitap bilgisini al
	book, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// Yazar bilgisini al
	authorInfo, err := s.authorService.GetAuthorInfo(book.Author)
	if err != nil {
		log.Printf("Yazar bilgisi alınamadı: %v", err)
		// Hata olsa bile kitap bilgisini döndür
		authorInfo = &model.AuthorInfo{
			Name:      book.Author,
			Biography: "Yazar bilgisi şu anda mevcut değil",
		}
	}

	return book.ToEnriched(authorInfo), nil
}

// GetBooksByAuthor yazar adına göre kitapları getirir
func (s *BookServiceImpl) GetBooksByAuthor(authorName string) ([]model.Book, error) {
	if authorName == "" {
		return nil, model.ErrInvalidAuthor
	}

	return s.bookRepo.GetBooksByAuthor(authorName)
}

// GetBooksByCategory kategori adına göre kitapları getirir
func (s *BookServiceImpl) GetBooksByCategory(categoryName string) ([]model.Book, error) {
	if categoryName == "" {
		return nil, errors.New("kategori adı boş olamaz")
	}

	return s.bookRepo.GetBooksByCategory(categoryName)
}

// GetBooksByCategoryWithPagination kategori adına göre sayfalanmış kitapları getirir
func (s *BookServiceImpl) GetBooksByCategoryWithPagination(categoryName string, params *model.BookSearchParams) (*model.PaginatedBooks, error) {
	if categoryName == "" {
		return nil, errors.New("kategori adı boş olamaz")
	}

	// Parametreleri doğrula
	if err := s.validateSearchParams(params); err != nil {
		return nil, err
	}

	// Kategori filtresini ayarla
	params.Category = categoryName

	return s.bookRepo.GetPaginatedBooks(params)
}

// GetEnrichedBooks yazar bilgisiyle zenginleştirilmiş kitap listesi getirir
func (s *BookServiceImpl) GetEnrichedBooks(params *model.BookSearchParams) ([]*model.EnrichedBook, error) {
	// Parametreleri doğrula
	if err := s.validateSearchParams(params); err != nil {
		return nil, err
	}

	// Zenginleştirilmiş kitaplar için sayfa boyutunu küçült
	if params.PageSize > 20 {
		params.PageSize = 10
	}

	result, err := s.bookRepo.GetPaginatedBooks(params)
	if err != nil {
		return nil, err
	}

	// Her kitap için yazar bilgisini zenginleştir
	enrichedBooks := make([]*model.EnrichedBook, len(result.Books))
	for i, book := range result.Books {
		authorInfo, err := s.authorService.GetAuthorInfo(book.Author)
		if err != nil {
			log.Printf("Yazar bilgisi alınamadı (%s): %v", book.Author, err)
			authorInfo = &model.AuthorInfo{
				Name:      book.Author,
				Biography: "Yazar bilgisi şu anda mevcut değil",
			}
		}

		enrichedBooks[i] = book.ToEnriched(authorInfo)
	}

	return enrichedBooks, nil
}

// validateSearchParams arama parametrelerini doğrular
func (s *BookServiceImpl) validateSearchParams(params *model.BookSearchParams) error {
	if params.Page < 1 {
		return model.ErrInvalidPage
	}
	
	if params.PageSize < 1 || params.PageSize > 100 {
		return model.ErrInvalidPageSize
	}

	return nil
} 