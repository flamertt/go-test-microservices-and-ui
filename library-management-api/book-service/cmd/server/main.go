package main

import (
	"database/sql"
	"log"

	"book-service/configs"
	"book-service/internal/handler"
	"book-service/internal/repository"
	"book-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Konfigürasyonu yükle
	cfg := configs.LoadConfig()

	// Veritabanı bağlantısını oluştur
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal("Veritabanı bağlantısı açılamadı:", err)
	}
	defer db.Close()

	// Bağlantıyı test et
	if err := db.Ping(); err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	log.Println("PostgreSQL veritabanına başarıyla bağlandı")

	// Dependency Injection -    katmanlarını oluştur
	bookRepo := repository.NewPostgreSQLBookRepository(db)
	authorService := service.NewHTTPAuthorService(cfg.Services.AuthorServiceURL)
	bookService := service.NewBookService(bookRepo, authorService)
	bookHandler := handler.NewBookHandler(bookService)

	// Gin router'ını oluştur
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "book-service",
			"message": "Book service çalışıyor",
		})
	})

	// API endpoint'leri - diğer servislerle tutarlılık için /api prefix'i kullan
	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/books", bookHandler.GetBooks)
		apiRoutes.GET("/books/:id", bookHandler.GetEnrichedBookByID) // Default olarak enriched döner
		apiRoutes.GET("/books/simple/:id", bookHandler.GetBookByID)  // Sadece kitap bilgisi
		apiRoutes.GET("/books/author/:authorName", bookHandler.GetBooksByAuthor)
		apiRoutes.GET("/books/category/:categoryName", bookHandler.GetBooksByCategory)
		apiRoutes.GET("/books/enriched", bookHandler.GetEnrichedBooks)
	}

	// Servisi başlat
	serverAddr := cfg.GetServerAddress()
	log.Printf("Book service %s adresinde başlatılıyor...", serverAddr)
	log.Println("🔗    Endpoints:")
	log.Println("  📚 GET /api/books                     - Sayfalı kitap listesi")
	log.Println("  📚 GET /api/books/:id                 - Zenginleştirilmiş kitap (yazar bilgisi ile)")
	log.Println("  📚 GET /api/books/simple/:id          - Sadece kitap bilgisi")
	log.Println("  📚 GET /api/books/author/:authorName  - Yazar kitapları")
	log.Println("  📚 GET /api/books/category/:category  - Kategori kitapları")
	log.Println("  📚 GET /api/books/enriched            - Zenginleştirilmiş kitap listesi")
	log.Println("  🩺 GET /health                       - Health check")
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
} 