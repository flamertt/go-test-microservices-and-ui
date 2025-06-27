package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"recommendation-service/configs"
	"recommendation-service/internal/handler"
	"recommendation-service/internal/service"
	"recommendation-service/pkg/logger"
)

func main() {
	// Load configuration
	cfg := configs.Load()

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize services
	bookService := service.NewBookService(cfg.BookServiceURL)
	authorService := service.NewAuthorService(cfg.AuthorServiceURL)
	genreService := service.NewGenreService(cfg.GenreServiceURL)
	recommendationService := service.NewRecommendationService(bookService, authorService, genreService)

	// Initialize handlers
	h := handler.NewRecommendationHandler(recommendationService, logger)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Health check
	router.GET("/health", h.HealthCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/recommendations", h.GetRecommendations)
		v1.GET("/recommendations/by-category", h.GetRecommendationsByCategory)
		v1.GET("/recommendations/by-author", h.GetRecommendationsByAuthor)
		v1.GET("/recommendations/trending", h.GetTrendingRecommendations)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	logger.Info("Starting recommendation service on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 