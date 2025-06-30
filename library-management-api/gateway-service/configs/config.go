package configs

import (
	"fmt"
	"os"
)

// Config uygulama konfigürasyonu
type Config struct {
	Server   ServerConfig   `json:"server"`
	Services ServicesConfig `json:"services"`
}

// ServerConfig server konfigürasyonu
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// ServicesConfig mikroservis URL'leri
type ServicesConfig struct {
	BookServiceURL          string `json:"book_service_url"`
	AuthorServiceURL        string `json:"author_service_url"`
	GenreServiceURL         string `json:"genre_service_url"`
	RecommendationServiceURL string `json:"recommendation_service_url"`
	AuthServiceURL          string `json:"auth_service_url"`
}

// LoadConfig konfigürasyonu yükler
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3000"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Services: ServicesConfig{
			BookServiceURL:          getEnv("BOOK_SERVICE_URL", "http://localhost:3001"),
			AuthorServiceURL:        getEnv("AUTHOR_SERVICE_URL", "http://localhost:3002"),
			GenreServiceURL:         getEnv("GENRE_SERVICE_URL", "http://localhost:3003"),
			RecommendationServiceURL: getEnv("RECOMMENDATION_SERVICE_URL", "http://localhost:3004"),
			AuthServiceURL:          getEnv("AUTH_SERVICE_URL", "http://localhost:3005"),
		},
	}
}

// GetServerAddress server adresini oluşturur
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// getEnv environment variable'ı okur, yoksa default değer döner
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 