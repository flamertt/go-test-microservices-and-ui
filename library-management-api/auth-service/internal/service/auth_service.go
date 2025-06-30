package service

import (
	"fmt"

	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/utils"
)

// AuthService authentication service interface'i
type AuthService interface {
	Register(req *model.RegisterRequest) (*model.RegisterResponse, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	ValidateToken(token string) (*model.JWTClaims, error)
	GetUserByID(userID uint) (*model.User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
	RefreshToken(token string) (string, error)
}

// authService authentication service implementasyonu
type authService struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewAuthService yeni auth service oluşturur
func NewAuthService(userRepo repository.UserRepository, jwtManager *utils.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register yeni kullanıcı kaydı yapar
func (s *authService) Register(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	// Şifre gücü kontrolü
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return nil, fmt.Errorf("şifre geçersiz: %w", err)
	}

	// Kullanıcı adı kontrolü
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı adı kontrolü yapılamadı: %w", err)
	}
	if exists {
		return nil, model.ErrUserAlreadyExists
	}

	// E-posta kontrolü
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("e-posta kontrolü yapılamadı: %w", err)
	}
	if exists {
		return nil, model.ErrUserAlreadyExists
	}

	// Şifreyi hash'le
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("şifre hash'lenemedi: %w", err)
	}

	// Yeni kullanıcı oluştur
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("kullanıcı oluşturulamadı: %w", err)
	}

	return &model.RegisterResponse{
		Message: "Kullanıcı başarıyla kaydedildi",
		User:    user.ToResponse(),
	}, nil
}

// Login kullanıcı girişi yapar
func (s *authService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Kullanıcıyı bul
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == model.ErrUserNotFound {
			return nil, model.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	// Şifreyi kontrol et
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, model.ErrInvalidCredentials
	}

	// JWT token oluştur
	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("token oluşturulamadı: %w", err)
	}

	return &model.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// ValidateToken token'ı doğrular
func (s *authService) ValidateToken(token string) (*model.JWTClaims, error) {
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token doğrulanamadı: %w", err)
	}

	// Kullanıcının hala mevcut olup olmadığını kontrol et
	_, err = s.userRepo.GetByID(claims.UserID)
	if err != nil {
		if err == model.ErrUserNotFound {
			return nil, model.ErrInvalidToken
		}
		return nil, fmt.Errorf("kullanıcı kontrol edilemedi: %w", err)
	}

	return claims, nil
}

// GetUserByID kullanıcıyı ID'ye göre getirir
func (s *authService) GetUserByID(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", err)
	}

	return user, nil
}

// ChangePassword kullanıcının şifresini değiştirir
func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Mevcut kullanıcıyı getir
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	// Eski şifreyi kontrol et
	if err := utils.CheckPassword(oldPassword, user.Password); err != nil {
		return model.ErrInvalidCredentials
	}

	// Yeni şifre gücünü kontrol et
	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("yeni şifre geçersiz: %w", err)
	}

	// Yeni şifreyi hash'le
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("şifre hash'lenemedi: %w", err)
	}

	// Şifreyi güncelle
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("şifre güncellenemedi: %w", err)
	}

	return nil
}

// RefreshToken token'ı yeniler
func (s *authService) RefreshToken(token string) (string, error) {
	// Mevcut token'ı doğrula
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return "", fmt.Errorf("token doğrulanamadı: %w", err)
	}

	// Kullanıcıyı getir
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	// Yeni token oluştur
	newToken, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("yeni token oluşturulamadı: %w", err)
	}

	return newToken, nil
} 