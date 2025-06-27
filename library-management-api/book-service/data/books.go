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

// Book API'de kullanılacak veri yapısı
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

// PaginatedBooks sayfalı kitap response yapısı
type PaginatedBooks struct {
	Books      []Book `json:"books"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// PaginatedAuthors sayfalı yazar response yapısı
type PaginatedAuthors struct {
	Authors    []map[string]interface{} `json:"authors"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

// PaginatedCategories sayfalı kategori response yapısı
type PaginatedCategories struct {
	Categories []map[string]interface{} `json:"categories"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

// ToBook BookDB'yi Book'a dönüştürür
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

	log.Println("PostgreSQL veritabanına başarıyla bağlandı")
	return nil
}

// GetAllBooks tüm kitapları getirir (backward compatibility)
func GetAllBooks() ([]Book, error) {
	result, err := GetPaginatedBooks(1, 50, "", "", "")
	if err != nil {
		return nil, err
	}
	return result.Books, nil
}

// GetPaginatedBooks sayfalı kitap listesi getirir
func GetPaginatedBooks(page, pageSize int, searchTerm, category, author string) (*PaginatedBooks, error) {
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

	// WHERE şartlarını hazırla
	whereClause := ""
	args := []interface{}{}
	argCount := 0

	if searchTerm != "" {
		argCount++
		whereClause += fmt.Sprintf(" WHERE (LOWER(book_title) LIKE LOWER($%d) OR LOWER(book_author) LIKE LOWER($%d) OR LOWER(book_publisher) LIKE LOWER($%d))", argCount, argCount, argCount)
		args = append(args, "%"+searchTerm+"%")
	}

	if category != "" {
		if whereClause != "" {
			argCount++
			whereClause += fmt.Sprintf(" AND LOWER(book_category_name) LIKE LOWER($%d)", argCount)
		} else {
			argCount++
			whereClause += fmt.Sprintf(" WHERE LOWER(book_category_name) LIKE LOWER($%d)", argCount)
		}
		args = append(args, "%"+category+"%")
	}

	if author != "" {
		if whereClause != "" {
			argCount++
			whereClause += fmt.Sprintf(" AND LOWER(book_author) LIKE LOWER($%d)", argCount)
		} else {
			argCount++
			whereClause += fmt.Sprintf(" WHERE LOWER(book_author) LIKE LOWER($%d)", argCount)
		}
		args = append(args, "%"+author+"%")
	}

	// Toplam sayıyı al
	countQuery := `SELECT COUNT(*) FROM ` + tableName + whereClause
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
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
	FROM ` + tableName + whereClause + `
	ORDER BY book_title 
	LIMIT $` + fmt.Sprintf("%d", argCount+1) + ` OFFSET $` + fmt.Sprintf("%d", argCount+2)

	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("kitaplar sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var bookDB BookDB
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

	totalPages := (total + pageSize - 1) / pageSize

	return &PaginatedBooks{
		Books:      books,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetBookByID ID'ye göre kitap getirir
func GetBookByID(id int) (*Book, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}

	query := `SELECT 
		$1 as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM ` + tableName + ` 
	ORDER BY book_title 
	LIMIT 1 OFFSET $2`

	var bookDB BookDB
	err := db.QueryRow(query, id, id-1).Scan(
		&bookDB.ID,
		&bookDB.Title,
		&bookDB.Publisher,
		&bookDB.Author,
		&bookDB.CategoryName,
		&bookDB.ProductCode,
		&bookDB.PageCount,
		&bookDB.ReleasedYear,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kitap bulunamadı")
	}
	if err != nil {
		return nil, fmt.Errorf("kitap sorgulanamadı: %v", err)
	}

	book := bookDB.ToBook()
	return &book, nil
}

// GetBooksByAuthor yazara göre kitapları getirir
func GetBooksByAuthor(authorName string) ([]Book, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}

	query := `SELECT 
		ROW_NUMBER() OVER (ORDER BY book_title) as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM ` + tableName + ` 
	WHERE LOWER(book_author) LIKE LOWER($1)`

	rows, err := db.Query(query, "%"+authorName+"%")
	if err != nil {
		return nil, fmt.Errorf("yazar kitapları sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var bookDB BookDB
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

// GetBooksByCategory kategoriye göre kitapları getirir
func GetBooksByCategory(categoryName string) ([]Book, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}

	query := `SELECT 
		ROW_NUMBER() OVER (ORDER BY book_title) as id,
		book_title, 
		book_publisher, 
		book_author, 
		book_category_name, 
		book_productcode, 
		book_page_count, 
		book_released_year 
	FROM ` + tableName + ` 
	WHERE LOWER(book_category_name) LIKE LOWER($1)`

	rows, err := db.Query(query, "%"+categoryName+"%")
	if err != nil {
		return nil, fmt.Errorf("kategori kitapları sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var bookDB BookDB
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

// GetAllAuthors tüm yazarları getirir (backward compatibility)
func GetAllAuthors() ([]string, error) {
	result, err := GetPaginatedAuthors(1, 1000, "")
	if err != nil {
		return nil, err
	}
	
	var authors []string
	for _, author := range result.Authors {
		if name, ok := author["name"].(string); ok {
			authors = append(authors, name)
		}
	}
	return authors, nil
}

// GetPaginatedAuthors sayfalı yazar listesi getirir
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
		return nil, fmt.Errorf("toplam yazar sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT DISTINCT book_author 
			  FROM ` + tableName + whereClause + `
			  ORDER BY book_author 
			  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("yazarlar sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var authors []map[string]interface{}
	id := (page-1)*pageSize + 1
	for rows.Next() {
		var author string
		if err := rows.Scan(&author); err != nil {
			continue
		}
		authors = append(authors, map[string]interface{}{
			"id":   id,
			"name": author,
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

// GetAllCategories tüm kategorileri getirir (backward compatibility)
func GetAllCategories() ([]string, error) {
	result, err := GetPaginatedCategories(1, 1000, "")
	if err != nil {
		return nil, err
	}
	
	var categories []string
	for _, category := range result.Categories {
		if name, ok := category["name"].(string); ok {
			categories = append(categories, name)
		}
	}
	return categories, nil
}

// GetPaginatedCategories sayfalı kategori listesi getirir
func GetPaginatedCategories(page, pageSize int, searchTerm string) (*PaginatedCategories, error) {
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
		return nil, fmt.Errorf("toplam kategori sayısı sorgulanamadı: %v", err)
	}

	// Sayfalı veriyi al
	query := `SELECT DISTINCT book_category_name 
			  FROM ` + tableName + whereClause + `
			  ORDER BY book_category_name 
			  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("kategoriler sorgulanamadı: %v", err)
	}
	defer rows.Close()

	var categories []map[string]interface{}
	id := (page-1)*pageSize + 1
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			continue
		}
		categories = append(categories, map[string]interface{}{
			"id":   id,
			"name": category,
		})
		id++
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &PaginatedCategories{
		Categories: categories,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// CloseDB veritabanı bağlantısını kapatır
func CloseDB() {
	if db != nil {
		db.Close()
	}
} 