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
		// Dinamik service routing - herhangi bir path'i ilgili servise yönlendir
		api.Any("/books/*path", h.RouteToService)
		api.Any("/books", h.RouteToService)
		
		api.Any("/authors/*path", h.RouteToService)
		api.Any("/authors", h.RouteToService)
		
		api.Any("/genres/*path", h.RouteToService)
		api.Any("/genres", h.RouteToService)
		
		api.Any("/recommendations/*path", h.RouteToService)
		api.Any("/recommendations", h.RouteToService)
		
		api.Any("/auth/*path", h.RouteToService)
		api.Any("/auth", h.RouteToService)
		
		// Services health check
		api.GET("/health", h.ServicesHealthCheck)
	}
}

func startServer(r *gin.Engine, cfg *configs.Config) {
	serverAddr := cfg.GetServerAddress()
	log.Printf("Gateway service %s adresinde başlatılıyor...", serverAddr)
	
	printAPIInfo()
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
}

func printAPIInfo() {
	log.Println("🚀 Microservices Gateway Başlatıldı")
	log.Println("📡 Dinamik Routing Aktif:")
	log.Println("  📚 /api/books/*        -> Book Service")
	log.Println("  ✍️  /api/authors/*      -> Author Service")
	log.Println("  📖 /api/genres/*        -> Genre Service")
	log.Println("  🤖 /api/recommendations/* -> Recommendation Service")
	log.Println("  🔐 /api/auth/*          -> Auth Service")
	log.Println("  🩺 /api/health          -> Services Health Check")
	log.Println("  🩺 /health              -> Gateway Health Check")
	log.Println("")
	log.Println("🔗 Gateway URL: http://localhost:3000")
} 