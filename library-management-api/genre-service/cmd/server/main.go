package main

import (
	"database/sql"
	"log"

	"genre-service/configs"
	"genre-service/internal/handler"
	"genre-service/internal/repository"
	"genre-service/internal/service"

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

	log.Println("Genre servisi PostgreSQL veritabanına başarıyla bağlandı")

	// Dependency Injection -    katmanlarını oluştur
	genreRepo := repository.NewPostgreSQLGenreRepository(db)
	bookService := service.NewHTTPBookService(cfg.Services.BookServiceURL)
	genreService := service.NewGenreService(genreRepo, bookService)
	genreHandler := handler.NewGenreHandler(genreService)

	// Gin router'ını oluştur
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "genre-service",
			"message": "Genre service çalışıyor",
		})
	})

	// API endpoint'leri - diğer servislerle tutarlılık için /api prefix'i kullan
	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/genres", genreHandler.GetGenres)
		apiRoutes.GET("/genres/:id", genreHandler.GetGenreByID)
		apiRoutes.GET("/genres/search", genreHandler.SearchGenres)
		apiRoutes.GET("/genres/detail/:name", genreHandler.GetGenreDetailByName)
	}

	// Servisi başlat
	serverAddr := cfg.GetServerAddress()
	log.Printf("Genre service %s adresinde başlatılıyor...", serverAddr)
	log.Println("🔗    Endpoints:")
	log.Println("  📖 GET /api/genres                     - Sayfalı tür listesi")
	log.Println("  📖 GET /api/genres/:id                 - Zenginleştirilmiş tür (kitap bilgisi ile)")
	log.Println("  📖 GET /api/genres/search?name=...     - Tür arama")
	log.Println("  📖 GET /api/genres/detail/:name        - Tür detayı + kitapları")
	log.Println("  🩺 GET /health                        - Health check")
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
} 