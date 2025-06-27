package handler

import (
	"net/http"
	"strconv"

	"book-service/internal/model"
	"book-service/internal/service"

	"github.com/gin-gonic/gin"
)

// BookHandler HTTP handler'ları
type BookHandler struct {
	bookService service.BookService
}

// NewBookHandler yeni book handler oluşturur
func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// GetBooks sayfalı kitap listesi endpoint'i
func (h *BookHandler) GetBooks(c *gin.Context) {
	params, err := h.parseSearchParams(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	result, err := h.bookService.GetPaginatedBooks(params)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "GET_BOOKS_ERROR", "Kitaplar getirilemedi")
		return
	}

	h.respondSuccess(c, result)
}

// GetBookByID ID'ye göre kitap getirme endpoint'i
func (h *BookHandler) GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_ID", "Geçersiz ID formatı")
		return
	}

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		if err == model.ErrBookNotFound {
			h.respondError(c, http.StatusNotFound, "BOOK_NOT_FOUND", "Kitap bulunamadı")
			return
		}
		h.respondError(c, http.StatusInternalServerError, "GET_BOOK_ERROR", "Kitap getirilemedi")
		return
	}

	h.respondSuccess(c, book)
}

// GetEnrichedBookByID zenginleştirilmiş kitap getirme endpoint'i
func (h *BookHandler) GetEnrichedBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_ID", "Geçersiz ID formatı")
		return
	}

	enrichedBook, err := h.bookService.GetEnrichedBookByID(id)
	if err != nil {
		if err == model.ErrBookNotFound {
			h.respondError(c, http.StatusNotFound, "BOOK_NOT_FOUND", "Kitap bulunamadı")
			return
		}
		h.respondError(c, http.StatusInternalServerError, "GET_ENRICHED_BOOK_ERROR", "Zenginleştirilmiş kitap getirilemedi")
		return
	}

	h.respondSuccess(c, enrichedBook)
}

// GetBooksByAuthor yazar adına göre kitaplar endpoint'i
func (h *BookHandler) GetBooksByAuthor(c *gin.Context) {
	authorName := c.Param("authorName")
	if authorName == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_AUTHOR", "Yazar adı gerekli")
		return
	}

	books, err := h.bookService.GetBooksByAuthor(authorName)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "GET_AUTHOR_BOOKS_ERROR", "Yazar kitapları getirilemedi")
		return
	}

	h.respondSuccess(c, books)
}

// GetBooksByCategory kategori adına göre kitaplar endpoint'i
func (h *BookHandler) GetBooksByCategory(c *gin.Context) {
	categoryName := c.Param("categoryName")
	if categoryName == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_CATEGORY", "Kategori adı gerekli")
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

	// Kategori parametresini set et
	params.Category = categoryName

	books, err := h.bookService.GetBooksByCategoryWithPagination(categoryName, params)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "GET_CATEGORY_BOOKS_ERROR", "Kategori kitapları getirilemedi")
		return
	}

	h.respondSuccess(c, books)
}

// GetEnrichedBooks zenginleştirilmiş kitap listesi endpoint'i
func (h *BookHandler) GetEnrichedBooks(c *gin.Context) {
	params, err := h.parseSearchParams(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	enrichedBooks, err := h.bookService.GetEnrichedBooks(params)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "GET_ENRICHED_BOOKS_ERROR", "Zenginleştirilmiş kitaplar getirilemedi")
		return
	}

	h.respondSuccess(c, enrichedBooks)
}

// parseSearchParams query parametrelerini parse eder
func (h *BookHandler) parseSearchParams(c *gin.Context) (*model.BookSearchParams, error) {
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

	return &model.BookSearchParams{
		Page:       page,
		PageSize:   pageSize,
		SearchTerm: c.Query("search"),
		Category:   c.Query("category"),
		Author:     c.Query("author"),
	}, nil
}

// respondSuccess başarılı yanıt gönderir
func (h *BookHandler) respondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// respondError hata yanıtı gönderir
func (h *BookHandler) respondError(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    errorCode,
			"message": message,
		},
	})
} 