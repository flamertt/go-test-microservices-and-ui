package model

// Genre domain model - kitap türü entity'si
type Genre struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PaginatedGenres sayfalı tür response yapısı
type PaginatedGenres struct {
	Genres     []Genre `json:"genres"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}

// EnrichedGenre kitap bilgisiyle zenginleştirilmiş tür
type EnrichedGenre struct {
	Genre
	Books     []BookInfo `json:"books"`
	BookCount int        `json:"book_count"`
}

// BookInfo book service'den gelen kitap bilgisi
type BookInfo struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Publisher    string `json:"publisher"`
	Author       string `json:"author"`
	CategoryName string `json:"category_name"`
	ProductCode  string `json:"product_code"`
	PageCount    int    `json:"page_count"`
	ReleasedYear int    `json:"released_year"`
}

// GenreSearchParams tür arama parametreleri
type GenreSearchParams struct {
	Page       int
	PageSize   int
	SearchTerm string
}

// Validate tür verilerini doğrular
func (g *Genre) Validate() error {
	if g.Name == "" {
		return ErrInvalidGenreName
	}
	return nil
}

// ToEnriched Genre'u EnrichedGenre'a dönüştürür
func (g *Genre) ToEnriched(books []BookInfo) *EnrichedGenre {
	return &EnrichedGenre{
		Genre:     *g,
		Books:     books,
		BookCount: len(books),
	}
} 