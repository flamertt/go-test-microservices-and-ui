package main

import (
	"database/sql"
	"log"
	"time"

	"auth-service/configs"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/utils"

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

	log.Println("Auth servisi PostgreSQL veritabanına başarıyla bağlandı")

	// JWT token süresini parse et
	tokenDuration, err := time.ParseDuration(cfg.JWT.TokenDuration)
	if err != nil {
		log.Fatal("Geçersiz token süresi:", err)
	}

	// Dependency Injection - katmanlarını oluştur
	userRepo := repository.NewPostgreSQLUserRepository(db)
	jwtManager := utils.NewJWTManager(cfg.JWT.SecretKey, tokenDuration)
	authService := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Gin router'ını oluştur
	r := gin.Default()

	// CORS middleware ekle
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "auth-service",
			"message": "Auth service çalışıyor",
		})
	})

	// API routes - diğer servislerle tutarlılık için /api prefix'i kullan
	apiRoutes := r.Group("/api")
	{
		// Public routes (authentication gerektirmeyen)
		publicRoutes := apiRoutes.Group("/auth")
		{
			publicRoutes.POST("/register", authHandler.Register)
			publicRoutes.POST("/login", authHandler.Login)
			publicRoutes.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (authentication gerektiren)
		protectedRoutes := apiRoutes.Group("/auth")
		protectedRoutes.Use(authMiddleware.RequireAuth())
		{
			protectedRoutes.GET("/profile", authHandler.GetProfile)
			protectedRoutes.POST("/change-password", authHandler.ChangePassword)
			protectedRoutes.GET("/validate", authHandler.ValidateToken)
			protectedRoutes.GET("/users/:id", authHandler.GetUser)
		}
	}

	// Servisi başlat
	serverAddr := cfg.GetServerAddress()
	log.Printf("Auth service %s adresinde başlatılıyor...", serverAddr)
	log.Println("🔗 Endpoints:")
	log.Println("  📝 POST /api/auth/register           - Kullanıcı kaydı")
	log.Println("  🔐 POST /api/auth/login              - Kullanıcı girişi")
	log.Println("  🔄 POST /api/auth/refresh            - Token yenileme")
	log.Println("  👤 GET  /api/auth/profile            - Kullanıcı profili (Protected)")
	log.Println("  🔑 POST /api/auth/change-password    - Şifre değiştirme (Protected)")
	log.Println("  ✅ GET  /api/auth/validate           - Token doğrulama (Protected)")
	log.Println("  👥 GET  /api/auth/users/:id          - Kullanıcı bilgisi (Protected)")
	log.Println("  🩺 GET  /health                     - Health check")
	log.Printf("  🔑 JWT Secret: %s", cfg.JWT.SecretKey[:10]+"...")
	log.Printf("  ⏰ Token Duration: %s", cfg.JWT.TokenDuration)

	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
} 