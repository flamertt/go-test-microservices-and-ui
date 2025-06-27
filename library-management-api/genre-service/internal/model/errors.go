package model

import "errors"

// Domain Error'ları - Genre business logic ihlalleri
var (
	ErrGenreNotFound        = errors.New("tür bulunamadı")
	ErrInvalidGenreID       = errors.New("geçersiz tür ID'si")
	ErrInvalidGenreName     = errors.New("tür adı boş olamaz")
	ErrInvalidPage          = errors.New("sayfa numarası 1'den küçük olamaz")
	ErrInvalidPageSize      = errors.New("sayfa boyutu 1-100 arasında olmalıdır")
	ErrDatabaseConnection   = errors.New("veritabanı bağlantı hatası")
	ErrBookServiceDown      = errors.New("kitap servisi kullanılamıyor")
)

// GenreError özel tür hatası
type GenreError struct {
	Code    string
	Message string
	Cause   error
}

func (e *GenreError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// NewGenreError yeni tür hatası oluşturur
func NewGenreError(code, message string, cause error) *GenreError {
	return &GenreError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
} 