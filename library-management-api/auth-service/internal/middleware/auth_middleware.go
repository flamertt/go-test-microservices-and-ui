package middleware

import (
	"net/http"

	"auth-service/internal/service"
	"auth-service/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware authentication middleware'ı
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware yeni auth middleware oluşturur
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth authentication gerektiren middleware
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header'ını al
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header bulunamadı",
			})
			c.Abort()
			return
		}

		// Token'ı extract et
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Geçersiz authorization header formatı",
			})
			c.Abort()
			return
		}

		// Token'ı doğrula
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized", 
				"message": "Geçersiz veya süresi dolmuş token",
			})
			c.Abort()
			return
		}

		// Claims'leri context'e ekle
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// RequireAuthOptional opsiyonel authentication middleware
func (m *AuthMiddleware) RequireAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header'ını al
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Token yoksa devam et
			c.Next()
			return
		}

		// Token'ı extract et
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Geçersiz format, devam et
			c.Next()
			return
		}

		// Token'ı doğrula
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			// Geçersiz token, devam et
			c.Next()
			return
		}

		// Claims'leri context'e ekle
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// CORS middleware'ı
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GetUserID context'ten user ID'yi alır
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	
	id, ok := userID.(uint)
	return id, ok
}

// GetUsername context'ten username'i alır
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	
	name, ok := username.(string)
	return name, ok
}

// GetEmail context'ten email'i alır
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	
	mail, ok := email.(string)
	return mail, ok
} 