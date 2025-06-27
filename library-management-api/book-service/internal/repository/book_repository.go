package repository

import (
	"database/sql"
	"fmt"
	"log"

	"book-service/internal/model"
	_ "github.com/lib/pq"
)

// BookRepository kitap veri erişim interface'i
type BookRepository interface {
	GetPaginatedBooks(params *model.BookSearchParams) (*model.PaginatedBooks, error)
	GetBookByID(id int) (*model.Book, error)
	GetBooksByAuthor(authorName string) ([]model.Book, error)
	GetBooksByCategory(categoryName string) ([]model.Book, error)
	Close() error
}

// PostgreSQLBookRepository PostgreSQL implementasyonu
type PostgreSQLBookRepository struct {
	db *sql.DB
}

// NewPostgreSQLBookRepository yeni PostgreSQL repository oluşturur
func NewPostgreSQLBookRepository(db *sql.DB) BookRepository {
	return &PostgreSQLBookRepository{
		db: db,
	}
}

// GetPaginatedBooks sayfalı kitap listesi getirir
func (r *PostgreSQLBookRepository) GetPaginatedBooks(params *model.BookSearchParams) (*model.PaginatedBooks, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 50
	}

	offset := (params.Page - 1) * params.PageSize

	// WHERE şartlarını hazırla
	whereClause := ""
	args := []interface{}{}
	argCount := 0

	if params.SearchTerm != "" {
		argCount++
		whereClause += fmt.Sprintf(" WHERE (LOWER(book_title) LIKE LOWER($%d) OR LOWER(book_author) LIKE LOWER($%d) OR LOWER(book_publisher) LIKE LOWER($%d))", argCount, argCount, argCount)
		args = append(args, "%"+params.SearchTerm+"%")
	}

	if params.Category != "" {
		if whereClause != "" {
			argCount++
			whereClause += fmt.Sprintf(" AND LOWER(book_category_name) LIKE LOWER($%d)", argCount)
		} else {
			argCount++
			whereClause += fmt.Sprintf(" WHERE LOWER(book_category_name) LIKE LOWER($%d)", argCount)
		}
		args = append(args, "%"+params.Category+"%")
	}

	if params.Author != "" {
		if whereClause != "" {
			argCount++
			whereClause += fmt.Sprintf(" AND LOWER(book_author) LIKE LOWER($%d)", argCount)
		} else {
			argCount++
			whereClause += fmt.Sprintf(" WHERE LOWER(book_author) LIKE LOWER($%d)", argCount)
		}
		args = append(args, "%"+params.Author+"%")
	}

	// Toplam sayıyı al
	countQuery := `SELECT COUNT(*) FROM books` + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("toplam kitap sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT 
		ROW_NUMBER() OVER (ORDER BY book_title) as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM books` + whereClause + `
	ORDER BY book_title 
	LIMIT $` + fmt.Sprintf("%d", argCount+1) + ` OFFSET $` + fmt.Sprintf("%d", argCount+2)

	args = append(args, params.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("kitaplar sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var bookDB model.BookDB
		err := rows.Scan(
			&bookDB.ID,
			&bookDB.Title,
			&bookDB.Publisher,
			&bookDB.Author,
			&bookDB.CategoryName,
			&bookDB.ProductCode,
			&bookDB.PageCount,
			&bookDB.ReleasedYear,
		)
		if err != nil {
			log.Printf("Kitap verisi okunamadı: %v", err)
			continue
		}
		books = append(books, bookDB.ToBook())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("satır okuma hatası: %v", err)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &model.PaginatedBooks{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetBookByID ID'ye göre kitap getirir
func (r *PostgreSQLBookRepository) GetBookByID(id int) (*model.Book, error) {
	query := `SELECT 
		$1 as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM books 
	ORDER BY book_title 
	LIMIT 1 OFFSET $2`

	var bookDB model.BookDB
	err := r.db.QueryRow(query, id, id-1).Scan(
		&bookDB.ID,
		&bookDB.Title,
		&bookDB.Publisher,
		&bookDB.Author,
		&bookDB.CategoryName,
		&bookDB.ProductCode,
		&bookDB.PageCount,
		&bookDB.ReleasedYear,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrBookNotFound
		}
		return nil, fmt.Errorf("kitap sorgulanamadı: %v", err)
	}

	book := bookDB.ToBook()
	return &book, nil
}

// GetBooksByAuthor yazar adına göre kitapları getirir
func (r *PostgreSQLBookRepository) GetBooksByAuthor(authorName string) ([]model.Book, error) {
	query := `SELECT 
		ROW_NUMBER() OVER (ORDER BY book_title) as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM books 
	WHERE LOWER(book_author) LIKE LOWER($1)
	ORDER BY book_title`

	rows, err := r.db.Query(query, "%"+authorName+"%")
	if err != nil {
		return nil, fmt.Errorf("yazar kitapları sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var bookDB model.BookDB
		err := rows.Scan(
			&bookDB.ID,
			&bookDB.Title,
			&bookDB.Publisher,
			&bookDB.Author,
			&bookDB.CategoryName,
			&bookDB.ProductCode,
			&bookDB.PageCount,
			&bookDB.ReleasedYear,
		)
		if err != nil {
			log.Printf("Kitap verisi okunamadı: %v", err)
			continue
		}
		books = append(books, bookDB.ToBook())
	}

	return books, nil
}

// GetBooksByCategory kategori adına göre kitapları getirir
func (r *PostgreSQLBookRepository) GetBooksByCategory(categoryName string) ([]model.Book, error) {
	query := `SELECT 
		ROW_NUMBER() OVER (ORDER BY book_title) as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM books 
	WHERE LOWER(book_category_name) LIKE LOWER($1)
	ORDER BY book_title`

	rows, err := r.db.Query(query, "%"+categoryName+"%")
	if err != nil {
		return nil, fmt.Errorf("kategori kitapları sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var bookDB model.BookDB
		err := rows.Scan(
			&bookDB.ID,
			&bookDB.Title,
			&bookDB.Publisher,
			&bookDB.Author,
			&bookDB.CategoryName,
			&bookDB.ProductCode,
			&bookDB.PageCount,
			&bookDB.ReleasedYear,
		)
		if err != nil {
			log.Printf("Kitap verisi okunamadı: %v", err)
			continue
		}
		books = append(books, bookDB.ToBook())
	}

	return books, nil
}

// Close veritabanı bağlantısını kapatır
func (r *PostgreSQLBookRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
} 