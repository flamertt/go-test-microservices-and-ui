package configs

import (
	"os"
)

type Config struct {
	Port             string
	Environment      string
	LogLevel         string
	BookServiceURL   string
	AuthorServiceURL string
	GenreServiceURL  string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "3004"),
		Environment:      getEnv("ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		BookServiceURL:   getEnv("BOOK_SERVICE_URL", "http://localhost:3001"),
		AuthorServiceURL: getEnv("AUTHOR_SERVICE_URL", "http://localhost:3002"),
		GenreServiceURL:  getEnv("GENRE_SERVICE_URL", "http://localhost:3003"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 