package handler

import (
	"net/http"
	"strconv"

	"genre-service/internal/model"
	"genre-service/internal/service"

	"github.com/gin-gonic/gin"
)

// GenreHandler HTTP handler'ları
type GenreHandler struct {
	genreService service.GenreService
}

// NewGenreHandler yeni genre handler oluşturur
func NewGenreHandler(genreService service.GenreService) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

// GetGenres sayfalı tür listesi endpoint'i
func (h *GenreHandler) GetGenres(c *gin.Context) {
	params, err := h.parseSearchParams(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	result, err := h.genreService.GetPaginatedGenres(params)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "GET_GENRES_ERROR", "Türler getirilemedi")
		return
	}

	h.respondSuccess(c, gin.H{
		"genres":      result.Genres,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// GetGenreByID ID'ye göre tür getirme endpoint'i (mock implementation)
func (h *GenreHandler) GetGenreByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_ID", "Geçersiz ID formatı")
		return
	}

	// Mock tür bilgisi (gerçekte DB'den gelecek)
	genreName := "Tür " + idParam
	
	// Zenginleştirilmiş tür bilgisini al
	enrichedGenre, err := h.genreService.GetEnrichedGenreByName(genreName)
	if err != nil {
		if err == model.ErrGenreNotFound {
			h.respondError(c, http.StatusNotFound, "GENRE_NOT_FOUND", "Tür bulunamadı")
			return
		}
		h.respondError(c, http.StatusInternalServerError, "GET_GENRE_ERROR", "Tür getirilemedi")
		return
	}

	h.respondSuccess(c, gin.H{
		"id":          id,
		"name":        enrichedGenre.Name,
		"description": enrichedGenre.Description,
		"books":       enrichedGenre.Books,
		"book_count":  enrichedGenre.BookCount,
	})
}

// GetGenreDetailByName tür adına göre detaylı bilgi endpoint'i
func (h *GenreHandler) GetGenreDetailByName(c *gin.Context) {
	genreName := c.Param("name")
	if genreName == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_GENRE_NAME", "Tür adı gerekli")
		return
	}

	// Pagination parametrelerini parse et
	params, err := h.parseSearchParams(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	// Maksimum 50 kitap limit'i
	if params.PageSize > 50 {
		params.PageSize = 50
	}

	enrichedGenre, err := h.genreService.GetEnrichedGenreByNameWithPagination(genreName, params)
	if err != nil {
		if err == model.ErrGenreNotFound {
			h.respondError(c, http.StatusNotFound, "GENRE_NOT_FOUND", "Tür bulunamadı")
			return
		}
		h.respondError(c, http.StatusInternalServerError, "GET_GENRE_DETAIL_ERROR", "Tür detayı getirilemedi")
		return
	}

	h.respondSuccess(c, gin.H{
		"genre":      enrichedGenre.Genre,
		"books":      enrichedGenre.Books,
		"book_count": enrichedGenre.BookCount,
		"page":       params.Page,
		"page_size":  params.PageSize,
		"total_pages": (enrichedGenre.BookCount + params.PageSize - 1) / params.PageSize,
		"total":      enrichedGenre.BookCount,
	})
}

// SearchGenres tür arama endpoint'i
func (h *GenreHandler) SearchGenres(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		h.respondError(c, http.StatusBadRequest, "MISSING_NAME", "name parametresi gerekli")
		return
	}

	genres, err := h.genreService.GetGenreByName(name)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "SEARCH_GENRES_ERROR", "Tür arama başarısız")
		return
	}

	h.respondSuccess(c, genres)
}

// parseSearchParams query parametrelerini parse eder
func (h *GenreHandler) parseSearchParams(c *gin.Context) (*model.GenreSearchParams, error) {
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

	return &model.GenreSearchParams{
		Page:       page,
		PageSize:   pageSize,
		SearchTerm: c.Query("search"),
	}, nil
}

// respondSuccess başarılı yanıt gönderir
func (h *GenreHandler) respondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// respondError hata yanıtı gönderir
func (h *GenreHandler) respondError(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    errorCode,
			"message": message,
		},
	})
} 