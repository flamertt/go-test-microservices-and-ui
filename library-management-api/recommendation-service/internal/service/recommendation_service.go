package service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"recommendation-service/internal/model"
)

type RecommendationService struct {
	bookService   *BookService
	authorService *AuthorService
	genreService  *GenreService
}

func NewRecommendationService(bookService *BookService, authorService *AuthorService, genreService *GenreService) *RecommendationService {
	return &RecommendationService{
		bookService:   bookService,
		authorService: authorService,
		genreService:  genreService,
	}
}

func (s *RecommendationService) GetGeneralRecommendations(limit int) (*model.RecommendationResponse, error) {
	books, err := s.bookService.GetAllBooks(100)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	recommendations := s.generateRecommendations(books, limit, "genel")
	
	return &model.RecommendationResponse{
		Recommendations: recommendations,
		Total:           len(recommendations),
		Timestamp:       time.Now().Format(time.RFC3339),
	}, nil
}

func (s *RecommendationService) GetRecommendationsByCategory(limit int) (*model.RecommendationResponse, error) {
	genres, err := s.genreService.GetAllGenres(50)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch genres: %w", err)
	}

	if len(genres) == 0 {
		return nil, fmt.Errorf("no genres found")
	}

	// Random genre seç
	rand.Seed(time.Now().UnixNano())
	selectedGenre := genres[rand.Intn(len(genres))]

	books, err := s.bookService.GetBooksByCategory(selectedGenre.Name, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books by category: %w", err)
	}

	recommendations := s.generateRecommendations(books, limit, fmt.Sprintf("%s kategorisi", selectedGenre.Name))

	return &model.RecommendationResponse{
		Recommendations: recommendations,
		Total:           len(recommendations),
		Category:        selectedGenre.Name,
		Timestamp:       time.Now().Format(time.RFC3339),
	}, nil
}

func (s *RecommendationService) GetRecommendationsByAuthor(limit int) (*model.RecommendationResponse, error) {
	authors, err := s.authorService.GetAllAuthors(100)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authors: %w", err)
	}

	if len(authors) == 0 {
		return nil, fmt.Errorf("no authors found")
	}

	// Random author seç
	rand.Seed(time.Now().UnixNano())
	selectedAuthor := authors[rand.Intn(len(authors))]

	books, err := s.bookService.GetBooksByAuthor(selectedAuthor.Name, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books by author: %w", err)
	}

	recommendations := s.generateRecommendations(books, limit, fmt.Sprintf("%s yazarından", selectedAuthor.Name))

	return &model.RecommendationResponse{
		Recommendations: recommendations,
		Total:           len(recommendations),
		Author:          selectedAuthor.Name,
		Timestamp:       time.Now().Format(time.RFC3339),
	}, nil
}

func (s *RecommendationService) GetTrendingRecommendations(limit int) (*model.RecommendationResponse, error) {
	books, err := s.bookService.GetAllBooks(100)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	// Yeni kitapları filtrele (son 10 yıl)
	currentYear := time.Now().Year()
	var recentBooks []model.Book
	for _, book := range books {
		if book.ReleasedYear >= currentYear-10 {
			recentBooks = append(recentBooks, book)
		}
	}

	if len(recentBooks) == 0 {
		recentBooks = books // Eğer yeni kitap yoksa hepsini kullan
	}

	recommendations := s.generateRecommendations(recentBooks, limit, "trend")
	
	return &model.RecommendationResponse{
		Recommendations: recommendations,
		Total:           len(recommendations),
		Timestamp:       time.Now().Format(time.RFC3339),
	}, nil
}

func (s *RecommendationService) generateRecommendations(books []model.Book, limit int, reasonType string) []model.Recommendation {
	if len(books) == 0 {
		return []model.Recommendation{}
	}

	// Kitapları karıştır
	rand.Seed(time.Now().UnixNano())
	shuffled := make([]model.Book, len(books))
	copy(shuffled, books)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Limit kadar al
	if len(shuffled) > limit {
		shuffled = shuffled[:limit]
	}

	recommendations := make([]model.Recommendation, len(shuffled))
	reasons := s.getReasons(reasonType)

	for i, book := range shuffled {
		recommendations[i] = model.Recommendation{
			Book:   book,
			Reason: reasons[rand.Intn(len(reasons))],
			Score:  rand.Intn(41) + 60, // 60-100 arası score
		}
	}

	return recommendations
}

func (s *RecommendationService) getReasons(reasonType string) []string {
	switch strings.ToLower(reasonType) {
	case "trend":
		return []string{
			"Şu anda trend olan popüler kitap",
			"Sosyal medyada çok konuşulan eser",
			"Bu ay en çok okunan kitaplardan",
			"Yeni çıkan ve beğenilen kitap",
			"Güncel ve popüler eser",
		}
	default:
		return []string{
			"Size özel seçilmiş kitap önerisi",
			"Okuma zevkinize uygun kitap",
			"Bu kitabı beğenebilirsiniz",
			"Kaliteli ve değerli eser",
			"Okuma listenize ekleyebileceğiniz kitap",
			"Farklı bir perspektif sunuyor",
			"Bu kategoriden önerilen kitap",
			"Yazarın en iyi eserlerinden",
			"Klasik ve zamansız eser",
			"Okurken keyif alacağınız kitap",
		}
	}
} 