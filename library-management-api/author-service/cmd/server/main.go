package main

import (
	"database/sql"
	"log"

	"author-service/configs"
	"author-service/internal/handler"
	"author-service/internal/repository"
	"author-service/internal/service"

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

	log.Println("Author servisi PostgreSQL veritabanına başarıyla bağlandı")

	// Dependency Injection -    katmanlarını oluştur
	authorRepo := repository.NewPostgreSQLAuthorRepository(db)
	bookService := service.NewHTTPBookService(cfg.Services.BookServiceURL)
	authorService := service.NewAuthorService(authorRepo, bookService)
	authorHandler := handler.NewAuthorHandler(authorService)

	// Gin router'ını oluştur
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "author-service",
			"message": "Author service çalışıyor",
		})
	})

	// API endpoint'leri - diğer servislerle tutarlılık için /api prefix'i kullan
	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/authors", authorHandler.GetAuthors)
		apiRoutes.GET("/authors/:id", authorHandler.GetAuthorByID)
		apiRoutes.GET("/authors/search", authorHandler.SearchAuthors)
		apiRoutes.GET("/authors/detail/:name", authorHandler.GetAuthorDetailByName)
	}

	// Servisi başlat
	serverAddr := cfg.GetServerAddress()
	log.Printf("Author service %s adresinde başlatılıyor...", serverAddr)
	log.Println("🔗    Endpoints:")
	log.Println("  ✍️  GET /api/authors                    - Sayfalı yazar listesi")
	log.Println("  ✍️  GET /api/authors/:id                - Zenginleştirilmiş yazar (kitap bilgisi ile)")
	log.Println("  ✍️  GET /api/authors/search?name=...    - Yazar arama")
	log.Println("  ✍️  GET /api/authors/detail/:name       - Yazar detayı + kitapları")
	log.Println("  🩺 GET /health                         - Health check")
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
} 