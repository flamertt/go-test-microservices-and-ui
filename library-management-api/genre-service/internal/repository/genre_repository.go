package repository

import (
	"database/sql"
	"fmt"

	"genre-service/internal/model"
	_ "github.com/lib/pq"
)

// GenreRepository tür veri erişim interface'i
type GenreRepository interface {
	GetPaginatedGenres(params *model.GenreSearchParams) (*model.PaginatedGenres, error)
	GetGenreByName(name string) ([]model.Genre, error)
	Close() error
}

// PostgreSQLGenreRepository PostgreSQL implementasyonu
type PostgreSQLGenreRepository struct {
	db *sql.DB
}

// NewPostgreSQLGenreRepository yeni PostgreSQL repository oluşturur
func NewPostgreSQLGenreRepository(db *sql.DB) GenreRepository {
	return &PostgreSQLGenreRepository{
		db: db,
	}
}

// GetPaginatedGenres sayfalı tür listesi getirir
func (r *PostgreSQLGenreRepository) GetPaginatedGenres(params *model.GenreSearchParams) (*model.PaginatedGenres, error) {
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
		whereClause = " WHERE LOWER(book_category_name) LIKE LOWER($1)"
		args = append(args, "%"+params.SearchTerm+"%")
	}

	// Toplam sayıyı al
	countQuery := `SELECT COUNT(DISTINCT book_category_name) FROM books` + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("toplam genre sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT DISTINCT book_category_name 
			  FROM books` + whereClause + `
			  ORDER BY book_category_name 
			  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	
	args = append(args, params.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("genreler sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var genres []model.Genre
	id := (params.Page-1)*params.PageSize + 1
	for rows.Next() {
		var genreName string
		if err := rows.Scan(&genreName); err != nil {
			continue
		}
		genres = append(genres, model.Genre{
			ID:          id,
			Name:        genreName,
			Description: fmt.Sprintf("%s kategorisindeki kitaplar", genreName),
		})
		id++
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &model.PaginatedGenres{
		Genres:     genres,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetGenreByName isim ile genre arama
func (r *PostgreSQLGenreRepository) GetGenreByName(name string) ([]model.Genre, error) {
	query := `SELECT DISTINCT book_category_name 
			  FROM books
			  WHERE LOWER(book_category_name) LIKE LOWER($1)
			  ORDER BY book_category_name`
	
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("genre arama sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var genres []model.Genre
	id := 1
	for rows.Next() {
		var genreName string
		if err := rows.Scan(&genreName); err != nil {
			continue
		}
		genres = append(genres, model.Genre{
			ID:          id,
			Name:        genreName,
			Description: fmt.Sprintf("%s kategorisindeki kitaplar", genreName),
		})
		id++
	}

	return genres, nil
}

// Close veritabanı bağlantısını kapatır
func (r *PostgreSQLGenreRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
} 