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

// Genre kitap türü veri yapısı
type Genre struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PaginatedGenres sayfalı genre response yapısı
type PaginatedGenres struct {
	Genres     []Genre `json:"genres"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
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

	log.Println("Genre servisi PostgreSQL veritabanına başarıyla bağlandı")
	return nil
}

// GetPaginatedGenres sayfalı genre listesi getirir
func GetPaginatedGenres(page, pageSize int, searchTerm string) (*PaginatedGenres, error) {
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
		whereClause = " WHERE LOWER(book_category_name) LIKE LOWER($1)"
		args = append(args, "%"+searchTerm+"%")
	}

	// Toplam sayıyı al
	countQuery := `SELECT COUNT(DISTINCT book_category_name) FROM ` + tableName + whereClause
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("toplam genre sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT DISTINCT book_category_name 
			  FROM ` + tableName + whereClause + `
			  ORDER BY book_category_name 
			  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("genreler sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var genres []Genre
	id := (page-1)*pageSize + 1
	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			continue
		}
		genres = append(genres, Genre{
			ID:          id,
			Name:        categoryName,
			Description: fmt.Sprintf("%s kategorisindeki kitapları keşfedin", categoryName),
		})
		id++
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &PaginatedGenres{
		Genres:     genres,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetGenreByName isim ile genre arama
func GetGenreByName(name string) ([]Genre, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}

	query := `SELECT DISTINCT book_category_name 
			  FROM ` + tableName + `
			  WHERE LOWER(book_category_name) LIKE LOWER($1)
			  ORDER BY book_category_name`
	
	rows, err := db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("genre arama sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var genres []Genre
	id := 1
	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			continue
		}
		genres = append(genres, Genre{
			ID:          id,
			Name:        categoryName,
			Description: fmt.Sprintf("%s kategorisindeki kitapları keşfedin", categoryName),
		})
		id++
	}

	return genres, nil
}

// CloseDB veritabanı bağlantısını kapatır
func CloseDB() {
	if db != nil {
		db.Close()
	}
} 