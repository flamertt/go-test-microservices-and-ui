package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyService proxy iş mantığı interface'i
type ProxyService interface {
	ProxyRequest(c *gin.Context, targetURL, basePath string)
	CheckServiceHealth(serviceURL string) string
	CheckAllServicesHealth(services map[string]string) map[string]string
}

// ProxyServiceImpl ProxyService implementasyonu
type ProxyServiceImpl struct {
	httpClient *http.Client
}

// NewProxyService yeni proxy service oluşturur
func NewProxyService() ProxyService {
	return &ProxyServiceImpl{
		httpClient: &http.Client{},
	}
}

// ProxyRequest HTTP isteğini hedef servise yönlendirir
func (s *ProxyServiceImpl) ProxyRequest(c *gin.Context, targetURL, basePath string) {
	// Path'i düzenle
	originalPath := c.Request.URL.Path
	
	// /api prefix'ini kaldır ve hedef service path'ini ekle
	var targetPath string
	if strings.HasPrefix(originalPath, "/api/recommendations") {
		// Recommendation service için özel handling ve endpoint mapping
		recommendationPath := strings.TrimPrefix(originalPath, "/api/recommendations")
		
		if strings.HasPrefix(recommendationPath, "/author") {
			targetPath = "/api/v1/recommendations/by-author" + strings.TrimPrefix(recommendationPath, "/author")
		} else if strings.HasPrefix(recommendationPath, "/category") {
			targetPath = "/api/v1/recommendations/by-category" + strings.TrimPrefix(recommendationPath, "/category")
		} else {
			// Ana recommendations endpoint'i
			targetPath = "/api/v1/recommendations" + recommendationPath
		}
	} else {
		// Diğer servisler için normal handling
		targetPath = strings.TrimPrefix(originalPath, "/api")
	}
	
	// Query parametrelerini ekle
	query := c.Request.URL.RawQuery
	if query != "" {
		targetPath += "?" + query
	}

	// Hedef URL'yi oluştur
	fullTargetURL := targetURL + targetPath

	log.Printf("Proxy: %s -> %s", originalPath, fullTargetURL)

	// HTTP isteği oluştur
	req, err := http.NewRequest(c.Request.Method, fullTargetURL, c.Request.Body)
	if err != nil {
		log.Printf("Proxy hatası (istek oluşturma): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "PROXY_REQUEST_ERROR",
				"message": "Proxy isteği oluşturulamadı",
			},
		})
		return
	}

	// Header'ları kopyala
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// İsteği gönder
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("Proxy hatası (servis bağlantısı): %v", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "SERVICE_UNAVAILABLE",
				"message": fmt.Sprintf("Servis bağlantı hatası: %v", err),
			},
		})
		return
	}
	defer resp.Body.Close()

	// Yanıtı oku
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Proxy hatası (yanıt okuma): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "RESPONSE_READ_ERROR",
				"message": "Yanıt okunamadı",
			},
		})
		return
	}

	// Content-Type'ı ayarla
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}

	// Yanıtı gönder
	c.Data(resp.StatusCode, contentType, body)
}

// CheckServiceHealth tek servisin health durumunu kontrol eder
func (s *ProxyServiceImpl) CheckServiceHealth(serviceURL string) string {
	resp, err := s.httpClient.Get(serviceURL + "/health")
	if err != nil {
		return "OFFLINE"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "OK"
	}
	return "ERROR"
}

// CheckAllServicesHealth tüm servislerin health durumunu kontrol eder
func (s *ProxyServiceImpl) CheckAllServicesHealth(services map[string]string) map[string]string {
	results := make(map[string]string)
	for serviceName, serviceURL := range services {
		results[serviceName] = s.CheckServiceHealth(serviceURL)
	}
	return results
} 