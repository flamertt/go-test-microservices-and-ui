package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword şifreyi hash'ler
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("şifre hash'lenemedi: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPassword şifreyi doğrular
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength şifre gücünü kontrol eder
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("şifre en az 6 karakter olmalıdır")
	}
	
	// Daha gelişmiş şifre kontrolleri burada eklenebilir
	// - Büyük harf kontrolü
	// - Küçük harf kontrolü  
	// - Sayı kontrolü
	// - Özel karakter kontrolü
	
	return nil
} 