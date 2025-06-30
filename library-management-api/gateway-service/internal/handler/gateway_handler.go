package handler

import (
	"net/http"
	"strings"

	"gateway-service/configs"
	"gateway-service/internal/service"

	"github.com/gin-gonic/gin"
)

// GatewayHandler HTTP handler'ları
type GatewayHandler struct {
	proxyService service.ProxyService
	config       *configs.Config
}

// NewGatewayHandler yeni gateway handler oluşturur
func NewGatewayHandler(proxyService service.ProxyService, config *configs.Config) *GatewayHandler {
	return &GatewayHandler{
		proxyService: proxyService,
		config:       config,
	}
}

// HealthCheck gateway health check endpoint'i
func (h *GatewayHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"service": "gateway-service",
		"message": "Gateway çalışıyor",
	})
}

// ServicesHealthCheck tüm servislerin health durumunu kontrol eder
func (h *GatewayHandler) ServicesHealthCheck(c *gin.Context) {
	services := map[string]string{
		"book-service":          h.config.Services.BookServiceURL,
		"author-service":        h.config.Services.AuthorServiceURL,
		"genre-service":         h.config.Services.GenreServiceURL,
		"recommendation-service": h.config.Services.RecommendationServiceURL,
		"auth-service":          h.config.Services.AuthServiceURL,
	}

	serviceHealths := h.proxyService.CheckAllServicesHealth(services)

	c.JSON(http.StatusOK, gin.H{
		"gateway":  "OK",
		"services": serviceHealths,
	})
}

// RouteToService dinamik olarak istekleri doğru servise yönlendirir
func (h *GatewayHandler) RouteToService(c *gin.Context) {
	// URL path'ini analiz et
	path := c.Request.URL.Path
	
	// Service'i path'e göre belirle
	var targetURL string
	var serviceName string
	
	switch {
	case strings.HasPrefix(path, "/api/books"):
		targetURL = h.config.Services.BookServiceURL
		serviceName = "book-service"
	case strings.HasPrefix(path, "/api/authors"):
		targetURL = h.config.Services.AuthorServiceURL
		serviceName = "author-service"
	case strings.HasPrefix(path, "/api/genres"):
		targetURL = h.config.Services.GenreServiceURL
		serviceName = "genre-service"
	case strings.HasPrefix(path, "/api/recommendations"):
		targetURL = h.config.Services.RecommendationServiceURL
		serviceName = "recommendation-service"
	case strings.HasPrefix(path, "/api/auth"):
		targetURL = h.config.Services.AuthServiceURL
		serviceName = "auth-service"
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "SERVICE_NOT_FOUND",
				"message": "İlgili servis bulunamadı",
				"path":    path,
			},
		})
		return
	}

	// İsteği ilgili servise yönlendir
	h.proxyService.ProxyRequest(c, targetURL, serviceName)
}

// Eski metodları kaldırıyoruz, artık RouteToService kullanacağız
// ProxyToBookService, ProxyToAuthorService, vs. metodları gereksiz 