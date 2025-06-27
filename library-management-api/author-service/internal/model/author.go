package model

// Author domain model - yazar entity'si
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PaginatedAuthors sayfalı yazar response yapısı
type PaginatedAuthors struct {
	Authors    []Author `json:"authors"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

// EnrichedAuthor kitap bilgisiyle zenginleştirilmiş yazar
type EnrichedAuthor struct {
	Author
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

// AuthorSearchParams yazar arama parametreleri
type AuthorSearchParams struct {
	Page       int
	PageSize   int
	SearchTerm string
}

// Validate yazar verilerini doğrular
func (a *Author) Validate() error {
	if a.Name == "" {
		return ErrInvalidAuthorName
	}
	return nil
}

// ToEnriched Author'u EnrichedAuthor'a dönüştürür
func (a *Author) ToEnriched(books []BookInfo) *EnrichedAuthor {
	return &EnrichedAuthor{
		Author:    *a,
		Books:     books,
		BookCount: len(books),
	}
} 