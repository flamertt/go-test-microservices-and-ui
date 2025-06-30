package model

import (
	"time"
)

// User kullanıcı modeli
type User struct {
	ID        uint      `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"` // JSON'da gösterilmez
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest giriş isteği yapısı
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest kayıt isteği yapısı  
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse giriş yanıt yapısı
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterResponse kayıt yanıt yapısı
type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

// JWTClaims JWT token claims yapısı
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Validate kullanıcı verilerini doğrular
func (u *User) Validate() error {
	if u.Username == "" {
		return ErrInvalidUsername
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

// ToResponse şifre olmadan kullanıcı bilgilerini döner
func (u *User) ToResponse() User {
	return User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
} 