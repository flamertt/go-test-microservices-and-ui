package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// ProxyService proxy iş mantığı interface'i
type ProxyService interface {
	ProxyRequest(c *gin.Context, targetURL, serviceName string)
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
func (s *ProxyServiceImpl) ProxyRequest(c *gin.Context, targetURL, serviceName string) {
	// Orijinal path'i al
	originalPath := c.Request.URL.Path
	
	// Tüm servisler tutarlı şekilde /api prefix'i kullanıyor
	targetPath := originalPath
	
	// Query parametrelerini ekle
	if c.Request.URL.RawQuery != "" {
		targetPath += "?" + c.Request.URL.RawQuery
	}

	// Hedef URL'yi oluştur
	fullTargetURL := targetURL + targetPath

	log.Printf("🔄 [%s] %s %s -> %s", 
		serviceName, 
		c.Request.Method, 
		originalPath, 
		fullTargetURL)

	// HTTP isteği oluştur
	req, err := http.NewRequest(c.Request.Method, fullTargetURL, c.Request.Body)
	if err != nil {
		log.Printf("❌ [%s] İstek oluşturma hatası: %v", serviceName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "PROXY_REQUEST_ERROR",
				"message": "Proxy isteği oluşturulamadı",
				"service": serviceName,
			},
		})
		return
	}

	// Header'ları kopyala (önemli olanları)
	s.copyHeaders(c.Request.Header, req.Header)

	// İsteği gönder
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("❌ [%s] Servis bağlantı hatası: %v", serviceName, err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "SERVICE_UNAVAILABLE",
				"message": fmt.Sprintf("%s servisi kullanılamıyor", serviceName),
				"service": serviceName,
				"details": err.Error(),
			},
		})
		return
	}
	defer resp.Body.Close()

	// Yanıtı oku
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [%s] Yanıt okuma hatası: %v", serviceName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "RESPONSE_READ_ERROR",
				"message": "Servis yanıtı okunamadı",
				"service": serviceName,
			},
		})
		return
	}

	// Response header'larını kopyala
	s.copyResponseHeaders(resp.Header, c.Writer.Header())

	// Başarılı proxy logla
	log.Printf("✅ [%s] %s %s -> %d (%d bytes)", 
		serviceName, 
		c.Request.Method, 
		originalPath, 
		resp.StatusCode,
		len(body))

	// Yanıtı gönder
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// copyHeaders önemli header'ları kopyalar
func (s *ProxyServiceImpl) copyHeaders(src, dst http.Header) {
	// Kritik header'ları kopyala
	importantHeaders := []string{
		"Authorization",
		"Content-Type",
		"Accept",
		"User-Agent",
		"X-Forwarded-For",
		"X-Real-IP",
	}
	
	for _, header := range importantHeaders {
		if values := src[header]; len(values) > 0 {
			for _, value := range values {
				dst.Add(header, value)
			}
		}
	}
}

// copyResponseHeaders yanıt header'larını kopyalar
func (s *ProxyServiceImpl) copyResponseHeaders(src, dst http.Header) {
	// Response header'larını kopyala
	responseHeaders := []string{
		"Content-Type",
		"Cache-Control",
		"Expires",
		"Last-Modified",
		"ETag",
	}
	
	for _, header := range responseHeaders {
		if values := src[header]; len(values) > 0 {
			for _, value := range values {
				dst.Set(header, value)
			}
		}
	}
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
		status := s.CheckServiceHealth(serviceURL)
		results[serviceName] = status
		log.Printf("🩺 [%s] Health Check: %s", serviceName, status)
	}
	return results
} 