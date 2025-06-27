package main

import (
	"log"

	"gateway-service/configs"
	"gateway-service/internal/handler"
	"gateway-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Konfigürasyonu yükle
	cfg := configs.LoadConfig()

	// Dependency Injection - katmanlarını oluştur
	proxyService := service.NewProxyService()
	gatewayHandler := handler.NewGatewayHandler(proxyService, cfg)

	// Gin router'ını oluştur
	r := gin.Default()

	// CORS ayarları
	setupCORS(r)

	// Route'ları ayarla
	setupRoutes(r, gatewayHandler)

	// Servisi başlat
	startServer(r, cfg)
}

func setupCORS(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))
}

func setupRoutes(r *gin.Engine, h *handler.GatewayHandler) {
	// Gateway health check
	r.GET("/health", h.HealthCheck)

	// API grubu
	api := r.Group("/api")
	{
		// Books Service
		books := api.Group("/books")
		books.GET("", h.ProxyToBookService).GET("/:id", h.ProxyToBookService).GET("/simple/:id", h.ProxyToBookService).GET("/author/:authorName", h.ProxyToBookService).GET("/category/:categoryName", h.ProxyToBookService).GET("/enriched", h.ProxyToBookService)

		// Authors Service  
		authors := api.Group("/authors")
		authors.GET("", h.ProxyToAuthorService).GET("/:id", h.ProxyToAuthorService).GET("/search", h.ProxyToAuthorService).GET("/detail/:name", h.ProxyToAuthorService)

		// Genres Service
		genres := api.Group("/genres")
		genres.GET("", h.ProxyToGenreService).GET("/:id", h.ProxyToGenreService).GET("/search", h.ProxyToGenreService).GET("/detail/:name", h.ProxyToGenreService)

		// Recommendations Service
		recommendations := api.Group("/recommendations")
		recommendations.GET("", h.ProxyToRecommendationService).GET("/category", h.ProxyToRecommendationService).GET("/author", h.ProxyToRecommendationService).GET("/category/:category", h.ProxyToRecommendationService).GET("/author/:author", h.ProxyToRecommendationService).GET("/status", h.ProxyToRecommendationService)
		
		// Services health check
		api.GET("/health", h.ServicesHealthCheck)
	}
}

func startServer(r *gin.Engine, cfg *configs.Config) {
	serverAddr := cfg.GetServerAddress()
	log.Printf("Gateway service %s adresinde başlatılıyor...", serverAddr)
	
	printAPIEndpoints()
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
}

func printAPIEndpoints() {
	log.Println("🔗    Microservices Gateway:")
	log.Println("  📚 Books (Book Service):         http://localhost:3000/api/books")
	log.Println("  📚 Books Enriched:               http://localhost:3000/api/books/enriched")
	log.Println("  ✍️  Authors (Author Service):     http://localhost:3000/api/authors")
	log.Println("  ✍️  Authors Detail:               http://localhost:3000/api/authors/detail/:name")
	log.Println("  📖 Genres (Genre Service):       http://localhost:3000/api/genres")
	log.Println("  📖 Genres Detail:                http://localhost:3000/api/genres/detail/:name")
	log.Println("  🤖 Recommendations:              http://localhost:3000/api/recommendations")
	log.Println("  🤖 Recommendations by Category:  http://localhost:3000/api/recommendations/category/:category")
	log.Println("  🤖 Recommendations by Author:    http://localhost:3000/api/recommendations/author/:author")
	log.Println("  🩺 Health Check:                 http://localhost:3000/api/health")
	log.Println("  🩺 Gateway Health:               http://localhost:3000/health")
} 