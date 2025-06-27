package service

import (
	"log"

	"genre-service/internal/model"
	"genre-service/internal/repository"
)

// GenreService tür iş mantığı interface'i
type GenreService interface {
	GetPaginatedGenres(params *model.GenreSearchParams) (*model.PaginatedGenres, error)
	GetGenreByName(name string) ([]model.Genre, error)
	GetEnrichedGenreByName(name string) (*model.EnrichedGenre, error)
	GetEnrichedGenreByNameWithPagination(name string, params *model.GenreSearchParams) (*model.EnrichedGenre, error)
}

// GenreServiceImpl GenreService implementasyonu
type GenreServiceImpl struct {
	genreRepo   repository.GenreRepository
	bookService BookService
}

// NewGenreService yeni genre service oluşturur
func NewGenreService(genreRepo repository.GenreRepository, bookService BookService) GenreService {
	return &GenreServiceImpl{
		genreRepo:   genreRepo,
		bookService: bookService,
	}
}

// GetPaginatedGenres sayfalı tür listesi getirir
func (s *GenreServiceImpl) GetPaginatedGenres(params *model.GenreSearchParams) (*model.PaginatedGenres, error) {
	// Parametreleri doğrula
	if err := s.validateSearchParams(params); err != nil {
		return nil, err
	}

	return s.genreRepo.GetPaginatedGenres(params)
}

// GetGenreByName isim ile tür arama
func (s *GenreServiceImpl) GetGenreByName(name string) ([]model.Genre, error) {
	if name == "" {
		return nil, model.ErrInvalidGenreName
	}

	return s.genreRepo.GetGenreByName(name)
}

// GetEnrichedGenreByName kitap bilgisiyle zenginleştirilmiş tür getirir
func (s *GenreServiceImpl) GetEnrichedGenreByName(name string) (*model.EnrichedGenre, error) {
	// Önce tür bilgisini al
	genres, err := s.GetGenreByName(name)
	if err != nil {
		return nil, err
	}

	if len(genres) == 0 {
		return nil, model.ErrGenreNotFound
	}

	// İlk türü al (genelde tek olur)
	genre := genres[0]

	// Book service'den bu türün kitaplarını al
	books, err := s.bookService.GetBooksByCategory(name)
	if err != nil {
		log.Printf("Tür kitapları alınamadı: %v", err)
		// Hata olsa bile tür bilgisini döndür
		books = []model.BookInfo{}
	}

	return genre.ToEnriched(books), nil
}

// GetEnrichedGenreByNameWithPagination kitap bilgisiyle zenginleştirilmiş tür getirir (sayfalanmış)
func (s *GenreServiceImpl) GetEnrichedGenreByNameWithPagination(name string, params *model.GenreSearchParams) (*model.EnrichedGenre, error) {
	// Önce tür bilgisini al
	genres, err := s.GetGenreByName(name)
	if err != nil {
		return nil, err
	}

	if len(genres) == 0 {
		return nil, model.ErrGenreNotFound
	}

	// İlk türü al (genelde tek olur)
	genre := genres[0]

	// Book service'den bu türün kitaplarını sayfalanmış olarak al
	books, err := s.bookService.GetBooksByCategoryWithPagination(name, params.Page, params.PageSize)
	if err != nil {
		log.Printf("Tür kitapları alınamadı: %v", err)
		// Hata olsa bile tür bilgisini döndür
		books = []model.BookInfo{}
	}

	// Toplam kitap sayısını al
	totalBooks, err := s.bookService.GetBookCountByCategory(name)
	if err != nil {
		log.Printf("Tür kitap sayısı alınamadı: %v", err)
		totalBooks = len(books)
	}

	enrichedGenre := genre.ToEnriched(books)
	enrichedGenre.BookCount = totalBooks

	return enrichedGenre, nil
}

// validateSearchParams arama parametrelerini doğrular
func (s *GenreServiceImpl) validateSearchParams(params *model.GenreSearchParams) error {
	if params.Page < 1 {
		return model.ErrInvalidPage
	}
	
	if params.PageSize < 1 || params.PageSize > 100 {
		return model.ErrInvalidPageSize
	}

	return nil
} 