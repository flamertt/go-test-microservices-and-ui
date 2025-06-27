package configs

import (
	"fmt"
	"os"
)

// Config uygulama konfigürasyonu
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Services ServicesConfig `json:"services"`
}

// ServerConfig server konfigürasyonu
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// DatabaseConfig veritabanı konfigürasyonu
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// ServicesConfig harici servis konfigürasyonları
type ServicesConfig struct {
	AuthorServiceURL string `json:"author_service_url"`
}

// LoadConfig konfigürasyonu yükler
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3001"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "mertpeker"),
			Password: getEnv("DB_PASSWORD", "mert123"),
			DBName:   getEnv("DB_NAME", "mertpeker"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Services: ServicesConfig{
			AuthorServiceURL: getEnv("AUTHOR_SERVICE_URL", "http://localhost:3002"),
		},
	}
}

// GetDatabaseURL veritabanı bağlantı string'ini oluşturur
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
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