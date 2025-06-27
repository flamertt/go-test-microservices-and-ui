package model

import "errors"

// Domain Error'ları - Author business logic ihlalleri
var (
	ErrAuthorNotFound       = errors.New("yazar bulunamadı")
	ErrInvalidAuthorID      = errors.New("geçersiz yazar ID'si")
	ErrInvalidAuthorName    = errors.New("yazar adı boş olamaz")
	ErrInvalidPage          = errors.New("sayfa numarası 1'den küçük olamaz")
	ErrInvalidPageSize      = errors.New("sayfa boyutu 1-100 arasında olmalıdır")
	ErrDatabaseConnection   = errors.New("veritabanı bağlantı hatası")
	ErrBookServiceDown      = errors.New("kitap servisi kullanılamıyor")
)

// AuthorError özel yazar hatası
type AuthorError struct {
	Code    string
	Message string
	Cause   error
}

func (e *AuthorError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// NewAuthorError yeni yazar hatası oluşturur
func NewAuthorError(code, message string, cause error) *AuthorError {
	return &AuthorError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
} 