package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"recommendation-service/internal/model"
	"recommendation-service/internal/service"
	"recommendation-service/pkg/logger"
)

type RecommendationHandler struct {
	service *service.RecommendationService
	logger  *logger.Logger
}

func NewRecommendationHandler(service *service.RecommendationService, logger *logger.Logger) *RecommendationHandler {
	return &RecommendationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *RecommendationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: gin.H{
			"status":  "healthy",
			"service": "recommendation-service",
		},
		Message: "Service is running",
	})
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	limit := h.getLimit(c)

	recommendations, err := h.service.GetGeneralRecommendations(limit)
	if err != nil {
		h.logger.Error("Failed to get recommendations: " + err.Error())
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "Failed to get recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    recommendations,
		Message: "Recommendations retrieved successfully",
	})
}

func (h *RecommendationHandler) GetRecommendationsByCategory(c *gin.Context) {
	limit := h.getLimit(c)

	recommendations, err := h.service.GetRecommendationsByCategory(limit)
	if err != nil {
		h.logger.Error("Failed to get category recommendations: " + err.Error())
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "Failed to get category recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    recommendations,
		Message: "Category recommendations retrieved successfully",
	})
}

func (h *RecommendationHandler) GetRecommendationsByAuthor(c *gin.Context) {
	limit := h.getLimit(c)

	recommendations, err := h.service.GetRecommendationsByAuthor(limit)
	if err != nil {
		h.logger.Error("Failed to get author recommendations: " + err.Error())
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "Failed to get author recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    recommendations,
		Message: "Author recommendations retrieved successfully",
	})
}

func (h *RecommendationHandler) GetTrendingRecommendations(c *gin.Context) {
	limit := h.getLimit(c)

	recommendations, err := h.service.GetTrendingRecommendations(limit)
	if err != nil {
		h.logger.Error("Failed to get trending recommendations: " + err.Error())
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "Failed to get trending recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    recommendations,
		Message: "Trending recommendations retrieved successfully",
	})
}

func (h *RecommendationHandler) getLimit(c *gin.Context) int {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}
	return limit
} 