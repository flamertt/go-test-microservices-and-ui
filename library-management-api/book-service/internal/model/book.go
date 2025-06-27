package model

import "database/sql"

// Book domain model - iş mantığının merkezindeki kitap entity'si
type Book struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Publisher    string `json:"publisher"`
	Author       string `json:"author"`
	CategoryName string `json:"category_name"`
	ProductCode  string `json:"product_code"`
	PageCount    int    `json:"page_count"`
	ReleasedYear int    `json:"released_year"`
}

// BookDB veritabanından gelen ham veri yapısı (NULL değerlerle)
type BookDB struct {
	ID           int            
	Title        sql.NullString         
	Publisher    sql.NullString         
	Author       sql.NullString         
	CategoryName sql.NullString         
	ProductCode  sql.NullString         
	PageCount    sql.NullInt32  
	ReleasedYear sql.NullInt32  
}

// PaginatedBooks sayfalı kitap response yapısı
type PaginatedBooks struct {
	Books      []Book `json:"books"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// EnrichedBook yazar bilgisiyle zenginleştirilmiş kitap
type EnrichedBook struct {
	Book
	AuthorInfo *AuthorInfo `json:"author_info"`
}

// AuthorInfo author service'den gelen yazar bilgisi
type AuthorInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Biography string `json:"biography"`
}

// BookSearchParams kitap arama parametreleri
type BookSearchParams struct {
	Page       int
	PageSize   int
	SearchTerm string
	Category   string
	Author     string
}

// ToBook BookDB'yi Book domain model'e dönüştürür
func (db BookDB) ToBook() Book {
	return Book{
		ID:           db.ID,
		Title:        db.Title.String,
		Publisher:    db.Publisher.String,
		Author:       db.Author.String,
		CategoryName: db.CategoryName.String,
		ProductCode:  db.ProductCode.String,
		PageCount:    int(db.PageCount.Int32),
		ReleasedYear: int(db.ReleasedYear.Int32),
	}
}

// Validate kitap verilerini doğrular
func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrInvalidTitle
	}
	if b.Author == "" {
		return ErrInvalidAuthor
	}
	return nil
}

// ToEnriched Book'u EnrichedBook'a dönüştürür
func (b *Book) ToEnriched(authorInfo *AuthorInfo) *EnrichedBook {
	return &EnrichedBook{
		Book:       *b,
		AuthorInfo: authorInfo,
	}
} 