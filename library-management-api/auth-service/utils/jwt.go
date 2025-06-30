package utils

import (
	"fmt"
	"time"

	"auth-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// JWTManager JWT token yöneticisi
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager yeni JWT manager oluşturur
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken kullanıcı için JWT token oluşturur
func (j *JWTManager) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(j.tokenDuration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("token oluşturulamadı: %w", err)
	}

	return tokenString, nil
}

// VerifyToken JWT token'ı doğrular ve claims'leri döner
func (j *JWTManager) VerifyToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method kontrolü
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmeyen signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse hatası: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("geçersiz token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims parse hatası")
	}

	// Claims'leri model'e dönüştür
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("geçersiz user_id")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("geçersiz username")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("geçersiz email")
	}

	return &model.JWTClaims{
		UserID:   uint(userID),
		Username: username,
		Email:    email,
	}, nil
}

// ExtractTokenFromHeader Authorization header'ından token'ı çıkarır
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header bulunamadı")
	}

	// "Bearer " prefix kontrolü
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("geçersiz authorization header formatı")
	}

	return authHeader[len(bearerPrefix):], nil
} 