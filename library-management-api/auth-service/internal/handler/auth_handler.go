package handler

import (
	"net/http"
	"strconv"

	"auth-service/internal/middleware"
	"auth-service/internal/model"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler authentication handler'ı
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler yeni auth handler oluşturur
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register kullanıcı kayıt endpoint'i
// @Summary Kullanıcı kaydı
// @Description Yeni kullanıcı kaydı yapar
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Kayıt bilgileri"
// @Success 201 {object} model.RegisterResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Geçersiz istek formatı: " + err.Error(),
		})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		if err == model.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": "Kullanıcı adı veya e-posta zaten kullanımda",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Kullanıcı kaydı yapılamadı: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login kullanıcı giriş endpoint'i
// @Summary Kullanıcı girişi
// @Description Kullanıcı girişi yapar ve JWT token döner
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Giriş bilgileri"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Geçersiz istek formatı: " + err.Error(),
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		if err == model.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Geçersiz kullanıcı adı veya şifre",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Giriş yapılamadı: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile kullanıcı profil bilgilerini getirir
// @Summary Kullanıcı profili
// @Description Giriş yapmış kullanıcının profil bilgilerini getirir
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Kullanıcı bilgisi bulunamadı",
		})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if err == model.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Kullanıcı bulunamadı",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Kullanıcı bilgisi getirilemedi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// ChangePassword kullanıcının şifresini değiştirir
// @Summary Şifre değiştir
// @Description Giriş yapmış kullanıcının şifresini değiştirir
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{old_password=string,new_password=string} true "Şifre değiştirme bilgileri"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Kullanıcı bilgisi bulunamadı",
		})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Geçersiz istek formatı: " + err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err == model.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Mevcut şifre yanlış",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Şifre değiştirilemedi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Şifre başarıyla değiştirildi",
	})
}

// RefreshToken token'ı yeniler
// @Summary Token yenile
// @Description Mevcut token ile yeni token alır
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{token=string} true "Yenilenecek token"
// @Success 200 {object} object{token=string}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Geçersiz istek formatı: " + err.Error(),
		})
		return
	}

	newToken, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Token yenilenemedi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// ValidateToken token doğrulama endpoint'i
// @Summary Token doğrula
// @Description Token'ın geçerli olup olmadığını kontrol eder
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{valid=bool,user_id=uint,username=string,email=string}
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Token geçersiz",
		})
		return
	}

	username, _ := middleware.GetUsername(c)
	email, _ := middleware.GetEmail(c)

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"user_id":  userID,
		"username": username,
		"email":    email,
	})
}

// GetUser ID'ye göre kullanıcı bilgilerini getirir
// @Summary Kullanıcı bilgisi
// @Description ID'ye göre kullanıcı bilgilerini getirir
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param id path int true "Kullanıcı ID"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/users/{id} [get]
func (h *AuthHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Geçersiz kullanıcı ID",
		})
		return
	}

	user, err := h.authService.GetUserByID(uint(id))
	if err != nil {
		if err == model.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Kullanıcı bulunamadı",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Kullanıcı bilgisi getirilemedi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
} 