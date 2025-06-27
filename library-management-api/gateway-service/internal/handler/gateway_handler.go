package handler

import (
	"net/http"

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
	}

	serviceHealths := h.proxyService.CheckAllServicesHealth(services)

	c.JSON(http.StatusOK, gin.H{
		"gateway":  "OK",
		"services": serviceHealths,
	})
}

// ProxyToBookService kitap servisine proxy
func (h *GatewayHandler) ProxyToBookService(c *gin.Context) {
	h.proxyService.ProxyRequest(c, h.config.Services.BookServiceURL, "/books")
}

// ProxyToAuthorService yazar servisine proxy
func (h *GatewayHandler) ProxyToAuthorService(c *gin.Context) {
	h.proxyService.ProxyRequest(c, h.config.Services.AuthorServiceURL, "/authors")
}

// ProxyToGenreService tür servisine proxy
func (h *GatewayHandler) ProxyToGenreService(c *gin.Context) {
	h.proxyService.ProxyRequest(c, h.config.Services.GenreServiceURL, "/genres")
}

// ProxyToRecommendationService öneri servisine proxy
func (h *GatewayHandler) ProxyToRecommendationService(c *gin.Context) {
	h.proxyService.ProxyRequest(c, h.config.Services.RecommendationServiceURL, "/recommendations")
} 