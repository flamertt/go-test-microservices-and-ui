package repository

import (
	"database/sql"
	"fmt"

	"author-service/internal/model"
	_ "github.com/lib/pq"
)

// AuthorRepository yazar veri erişim interface'i
type AuthorRepository interface {
	GetPaginatedAuthors(params *model.AuthorSearchParams) (*model.PaginatedAuthors, error)
	GetAuthorByName(name string) ([]model.Author, error)
	Close() error
}

// PostgreSQLAuthorRepository PostgreSQL implementasyonu
type PostgreSQLAuthorRepository struct {
	db *sql.DB
}

// NewPostgreSQLAuthorRepository yeni PostgreSQL repository oluşturur
func NewPostgreSQLAuthorRepository(db *sql.DB) AuthorRepository {
	return &PostgreSQLAuthorRepository{
		db: db,
	}
}

// GetPaginatedAuthors sayfalı yazar listesi getirir
func (r *PostgreSQLAuthorRepository) GetPaginatedAuthors(params *model.AuthorSearchParams) (*model.PaginatedAuthors, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 50
	}

	offset := (params.Page - 1) * params.PageSize

	// WHERE şartını hazırla
	whereClause := ""
	args := []interface{}{}
	
	if params.SearchTerm != "" {
		whereClause = " WHERE LOWER(book_author) LIKE LOWER($1)"
		args = append(args, "%"+params.SearchTerm+"%")
	}

	// Toplam sayıyı al
	countQuery := `SELECT COUNT(DISTINCT book_author) FROM books` + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("toplam author sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT DISTINCT book_author 
			  FROM books` + whereClause + `
			  ORDER BY book_author 
			  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	
	args = append(args, params.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("authorlar sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var authors []model.Author
	id := (params.Page-1)*params.PageSize + 1
	for rows.Next() {
		var authorName string
		if err := rows.Scan(&authorName); err != nil {
			continue
		}
		authors = append(authors, model.Author{
			ID:   id,
			Name: authorName,
		})
		id++
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &model.PaginatedAuthors{
		Authors:    authors,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAuthorByName isim ile author arama
func (r *PostgreSQLAuthorRepository) GetAuthorByName(name string) ([]model.Author, error) {
	query := `SELECT DISTINCT book_author 
			  FROM books
			  WHERE LOWER(book_author) LIKE LOWER($1)
			  ORDER BY book_author`
	
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("author arama sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var authors []model.Author
	id := 1
	for rows.Next() {
		var authorName string
		if err := rows.Scan(&authorName); err != nil {
			continue
		}
		authors = append(authors, model.Author{
			ID:   id,
			Name: authorName,
		})
		id++
	}

	return authors, nil
}

// Close veritabanı bağlantısını kapatır
func (r *PostgreSQLAuthorRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
} 