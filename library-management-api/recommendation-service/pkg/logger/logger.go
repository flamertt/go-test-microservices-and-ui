package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New(level string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[RECOMMENDATION] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string) {
	l.Logger.Printf("[INFO] %s", msg)
}

func (l *Logger) Error(msg string) {
	l.Logger.Printf("[ERROR] %s", msg)
}

func (l *Logger) Debug(msg string) {
	l.Logger.Printf("[DEBUG] %s", msg)
}

func (l *Logger) Warn(msg string) {
	l.Logger.Printf("[WARN] %s", msg)
} 