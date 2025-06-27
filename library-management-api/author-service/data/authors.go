package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// PostgreSQL veritabanı konfigürasyonu
const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "mertpeker"
	dbPassword = "mert123"
	dbName     = "mertpeker"
	tableName  = "books"
)

// Author yazar veri yapısı - sadece isim var
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PaginatedAuthors sayfalı author response yapısı
type PaginatedAuthors struct {
	Authors    []Author `json:"authors"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

var db *sql.DB

// InitDB veritabanı bağlantısını başlatır
func InitDB() error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("veritabanı bağlantısı açılamadı: %v", err)
	}

	// Bağlantıyı test et
	if err = db.Ping(); err != nil {
		return fmt.Errorf("veritabanına bağlanılamadı: %v", err)
	}

	log.Println("Author servisi PostgreSQL veritabanına başarıyla bağlandı")
	return nil
}

// GetPaginatedAuthors sayfalı author listesi getirir
func GetPaginatedAuthors(page, pageSize int, searchTerm string) (*PaginatedAuthors, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// WHERE şartını hazırla
	whereClause := ""
	args := []interface{}{}
	
	if searchTerm != "" {
		whereClause = " WHERE LOWER(book_author) LIKE LOWER($1)"
		args = append(args, "%"+searchTerm+"%")
	}

	// Toplam sayıyı al
	countQuery := `SELECT COUNT(DISTINCT book_author) FROM ` + tableName + whereClause
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("toplam author sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT DISTINCT book_author 
			  FROM ` + tableName + whereClause + `
			  ORDER BY book_author 
			  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("authorlar sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var authors []Author
	id := (page-1)*pageSize + 1
	for rows.Next() {
		var authorName string
		if err := rows.Scan(&authorName); err != nil {
			continue
		}
		authors = append(authors, Author{
			ID:   id,
			Name: authorName,
		})
		id++
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &PaginatedAuthors{
		Authors:    authors,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAuthorByName isim ile author arama
func GetAuthorByName(name string) ([]Author, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}

	query := `SELECT DISTINCT book_author 
			  FROM ` + tableName + `
			  WHERE LOWER(book_author) LIKE LOWER($1)
			  ORDER BY book_author`
	
	rows, err := db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("author arama sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var authors []Author
	id := 1
	for rows.Next() {
		var authorName string
		if err := rows.Scan(&authorName); err != nil {
			continue
		}
		authors = append(authors, Author{
			ID:   id,
			Name: authorName,
		})
		id++
	}

	return authors, nil
}

// CloseDB veritabanı bağlantısını kapatır
func CloseDB() {
	if db != nil {
		db.Close()
	}
} 