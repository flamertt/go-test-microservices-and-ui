package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config uygulama konfigürasyon yapısı
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// ServerConfig server konfigürasyonu
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Environment string `mapstructure:"environment"`
	AppName     string `mapstructure:"app_name"`
	Version     string `mapstructure:"version"`
}

// DatabaseConfig veritabanı konfigürasyonu
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// LoggerConfig logger konfigürasyonu
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	OutputPath string `mapstructure:"output_path"`
}

// MetricsConfig metrics konfigürasyonu
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// Load konfigürasyonu yükler
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/go-test-microservice/")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Varsayılan değerler
	setDefaults()

	// Konfigürasyon dosyasını oku
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Konfigürasyon dosyası okunamadı: %v", err)
		log.Println("Varsayılan değerler ve environment variables kullanılacak")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults varsayılan değerleri ayarlar
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "3000")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.app_name", "Go Test Mikroservis")
	viper.SetDefault("server.version", "1.0.0")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "microservice_db")
	viper.SetDefault("database.sslmode", "disable")

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.encoding", "json")
	viper.SetDefault("logger.output_path", "stdout")

	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
}

// GetString wrapper for viper.GetString
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt wrapper for viper.GetInt
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool wrapper for viper.GetBool
func GetBool(key string) bool {
	return viper.GetBool(key)
} 