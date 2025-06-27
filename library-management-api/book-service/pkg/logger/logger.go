package logger

import (
	"go-test-microservice/pkg/config"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Init logger'ı başlatır
func Init(cfg *config.Config) error {
	var zapConfig zap.Config

	// Environment'a göre config seç
	if cfg.Server.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// Log level ayarla
	level, err := zapcore.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zap.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Encoding ayarla
	if cfg.Logger.Encoding == "console" {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig.Encoding = "json"
	}

	// Encoder config
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.CallerKey = "caller"
	zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Output path ayarla
	if cfg.Logger.OutputPath != "" && cfg.Logger.OutputPath != "stdout" {
		zapConfig.OutputPaths = []string{cfg.Logger.OutputPath}
		zapConfig.ErrorOutputPaths = []string{cfg.Logger.OutputPath}
	}

	// Logger oluştur
	Logger, err = zapConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	} else {
		log.Println("INFO:", msg)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	} else {
		log.Println("ERROR:", msg)
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	} else {
		log.Println("DEBUG:", msg)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	} else {
		log.Println("WARN:", msg)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	} else {
		log.Fatal("FATAL:", msg)
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// GetLogger returns the current logger instance
func GetLogger() *zap.Logger {
	return Logger
}

// SetLevel dinamik olarak log seviyesini değiştirir
func SetLevel(level string) {
	if Logger == nil {
		return
	}

	_, err := zapcore.ParseLevel(level)
	if err != nil {
		Error("Geçersiz log seviyesi", zap.String("level", level), zap.Error(err))
		return
	}

	// Yeni level'ı ayarla (production config gerekebilir)
	Info("Log seviyesi değiştirildi", zap.String("new_level", level))
}

// WithFields returns a logger with additional fields
func WithFields(fields ...zap.Field) *zap.Logger {
	if Logger == nil {
		return nil
	}
	return Logger.With(fields...)
}

// LogHTTPRequest HTTP isteklerini loglar
func LogHTTPRequest(method, path, ip, userAgent string, statusCode int, duration string) {
	Info("HTTP İstek",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", ip),
		zap.String("user_agent", userAgent),
		zap.Int("status_code", statusCode),
		zap.String("duration", duration),
	)
}

// LogError hataları loglar
func LogError(err error, context string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.Error(err),
		zap.String("context", context),
	}, fields...)
	
	Error("Hata oluştu", allFields...)
}

// GetLogLevel mevcut log seviyesini döndürür
func GetLogLevel() string {
	if Logger == nil {
		return "unknown"
	}
	return Logger.Level().String()
}

// IsDebugEnabled debug seviyesinin aktif olup olmadığını kontrol eder
func IsDebugEnabled() bool {
	if Logger == nil {
		return false
	}
	return Logger.Level() <= zap.DebugLevel
}

// NewRequestLogger her istek için yeni bir logger oluşturur
func NewRequestLogger(requestID string) *zap.Logger {
	if Logger == nil {
		return nil
	}
	return Logger.With(zap.String("request_id", requestID))
}

// InitFromEnv environment variables'dan logger'ı başlatır
func InitFromEnv() error {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	encoding := os.Getenv("LOG_ENCODING")
	if encoding == "" {
		encoding = "json"
	}

	outputPath := os.Getenv("LOG_OUTPUT_PATH")
	if outputPath == "" {
		outputPath = "stdout"
	}

	// Basit config oluştur
	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Level:      level,
			Encoding:   encoding,
			OutputPath: outputPath,
		},
		Server: config.ServerConfig{
			Environment: os.Getenv("GO_ENV"),
		},
	}

	return Init(cfg)
} 