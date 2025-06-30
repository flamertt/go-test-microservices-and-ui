package model

import "errors"

// Auth servisi hata tanımları
var (
	// Kullanıcı hataları
	ErrUserNotFound       = errors.New("kullanıcı bulunamadı")
	ErrUserAlreadyExists  = errors.New("kullanıcı zaten mevcut")
	ErrInvalidCredentials = errors.New("geçersiz kullanıcı adı veya şifre")
	ErrInvalidUsername    = errors.New("geçersiz kullanıcı adı")
	ErrInvalidEmail       = errors.New("geçersiz e-posta adresi")
	ErrWeakPassword       = errors.New("şifre çok zayıf")
	
	// Token hataları
	ErrInvalidToken       = errors.New("geçersiz token")
	ErrExpiredToken       = errors.New("token süresi dolmuş")
	ErrMissingToken       = errors.New("token bulunamadı")
	ErrTokenGeneration    = errors.New("token oluşturulamadı")
	
	// Genel hatalar
	ErrInternalServer     = errors.New("sunucu hatası")
	ErrInvalidRequest     = errors.New("geçersiz istek")
	ErrUnauthorized       = errors.New("yetkisiz erişim")
	ErrForbidden          = errors.New("yasaklı erişim")
)

// ErrorResponse API hata yanıt yapısı
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
} 