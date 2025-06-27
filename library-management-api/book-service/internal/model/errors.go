package model

import "errors"

// Domain Error'ları - İş mantığı kuralları ihlalleri
var (
	ErrBookNotFound      = errors.New("kitap bulunamadı")
	ErrInvalidBookID     = errors.New("geçersiz kitap ID'si")
	ErrInvalidTitle      = errors.New("kitap başlığı boş olamaz")
	ErrInvalidAuthor     = errors.New("yazar adı boş olamaz")
	ErrInvalidPage       = errors.New("sayfa numarası 1'den küçük olamaz")
	ErrInvalidPageSize   = errors.New("sayfa boyutu 1-100 arasında olmalıdır")
	ErrDatabaseConnection = errors.New("veritabanı bağlantı hatası")
	ErrAuthorServiceDown  = errors.New("yazar servisi kullanılamıyor")
)

// BookError özel kitap hatası
type BookError struct {
	Code    string
	Message string
	Cause   error
}

func (e *BookError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// NewBookError yeni kitap hatası oluşturur
func NewBookError(code, message string, cause error) *BookError {
	return &BookError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
} 