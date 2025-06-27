package handler

import (
	"net/http"
	"strconv"

	"author-service/internal/model"
	"author-service/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthorHandler HTTP handler'ları
type AuthorHandler struct {
	authorService service.AuthorService
}

// NewAuthorHandler yeni author handler oluşturur
func NewAuthorHandler(authorService service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		authorService: authorService,
	}
}

// GetAuthors sayfalı yazar listesi endpoint'i
func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	params, err := h.parseSearchParams(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	result, err := h.authorService.GetPaginatedAuthors(params)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "GET_AUTHORS_ERROR", "Yazarlar getirilemedi")
		return
	}

	h.respondSuccess(c, gin.H{
		"authors":     result.Authors,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// GetAuthorByID ID'ye göre yazar getirme endpoint'i (mock implementation)
func (h *AuthorHandler) GetAuthorByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_ID", "Geçersiz ID formatı")
		return
	}

	// Mock yazar bilgisi (gerçekte DB'den gelecek)
	authorName := "Yazar " + idParam
	
	// Zenginleştirilmiş yazar bilgisini al
	enrichedAuthor, err := h.authorService.GetEnrichedAuthorByName(authorName)
	if err != nil {
		if err == model.ErrAuthorNotFound {
			h.respondError(c, http.StatusNotFound, "AUTHOR_NOT_FOUND", "Yazar bulunamadı")
			return
		}
		h.respondError(c, http.StatusInternalServerError, "GET_AUTHOR_ERROR", "Yazar getirilemedi")
		return
	}

	h.respondSuccess(c, gin.H{
		"id":         id,
		"name":       enrichedAuthor.Name,
		"biography":  "Bu yazar için detay bilgisi",
		"books":      enrichedAuthor.Books,
		"book_count": enrichedAuthor.BookCount,
	})
}

// GetAuthorDetailByName yazar adına göre detaylı bilgi endpoint'i
func (h *AuthorHandler) GetAuthorDetailByName(c *gin.Context) {
	authorName := c.Param("name")
	if authorName == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_AUTHOR_NAME", "Yazar adı gerekli")
		return
	}

	enrichedAuthor, err := h.authorService.GetEnrichedAuthorByName(authorName)
	if err != nil {
		if err == model.ErrAuthorNotFound {
			h.respondError(c, http.StatusNotFound, "AUTHOR_NOT_FOUND", "Yazar bulunamadı")
			return
		}
		h.respondError(c, http.StatusInternalServerError, "GET_AUTHOR_DETAIL_ERROR", "Yazar detayı getirilemedi")
		return
	}

	h.respondSuccess(c, gin.H{
		"author":     enrichedAuthor.Author,
		"books":      enrichedAuthor.Books,
		"book_count": enrichedAuthor.BookCount,
	})
}

// SearchAuthors yazar arama endpoint'i
func (h *AuthorHandler) SearchAuthors(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		h.respondError(c, http.StatusBadRequest, "MISSING_NAME", "name parametresi gerekli")
		return
	}

	authors, err := h.authorService.GetAuthorByName(name)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "SEARCH_AUTHORS_ERROR", "Yazar arama başarısız")
		return
	}

	h.respondSuccess(c, authors)
}

// parseSearchParams query parametrelerini parse eder
func (h *AuthorHandler) parseSearchParams(c *gin.Context) (*model.AuthorSearchParams, error) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	return &model.AuthorSearchParams{
		Page:       page,
		PageSize:   pageSize,
		SearchTerm: c.Query("search"),
	}, nil
}

// respondSuccess başarılı yanıt gönderir
func (h *AuthorHandler) respondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// respondError hata yanıtı gönderir
func (h *AuthorHandler) respondError(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    errorCode,
			"message": message,
		},
	})
} 