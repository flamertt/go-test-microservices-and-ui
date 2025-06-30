package repository

import (
	"database/sql"
	"fmt"
	"time"

	"auth-service/internal/model"
)

// UserRepository kullanıcı repository interface'i
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
}

// postgresUserRepository PostgreSQL user repository implementasyonu
type postgresUserRepository struct {
	db *sql.DB
}

// NewPostgreSQLUserRepository yeni PostgreSQL user repository oluşturur
func NewPostgreSQLUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

// Create yeni kullanıcı oluşturur
func (r *postgresUserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	
	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("kullanıcı oluşturulamadı: %w", err)
	}
	
	return nil
}

// GetByID ID'ye göre kullanıcı getirir
func (r *postgresUserRepository) GetByID(id uint) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = $1`
	
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", err)
	}
	
	return user, nil
}

// GetByUsername kullanıcı adına göre kullanıcı getirir
func (r *postgresUserRepository) GetByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE username = $1`
	
	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", err)
	}
	
	return user, nil
}

// GetByEmail e-postaya göre kullanıcı getirir
func (r *postgresUserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1`
	
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", err)
	}
	
	return user, nil
}

// Update kullanıcı bilgilerini günceller
func (r *postgresUserRepository) Update(user *model.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, password_hash = $4, updated_at = $5 
		WHERE id = $1`
	
	user.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("kullanıcı güncellenemedi: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("güncelleme sonucu alınamadı: %w", err)
	}
	
	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}
	
	return nil
}

// Delete kullanıcıyı siler
func (r *postgresUserRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("kullanıcı silinemedi: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("silme sonucu alınamadı: %w", err)
	}
	
	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}
	
	return nil
}

// ExistsByUsername kullanıcı adının mevcut olup olmadığını kontrol eder
func (r *postgresUserRepository) ExistsByUsername(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	
	var exists bool
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("kullanıcı adı kontrolü yapılamadı: %w", err)
	}
	
	return exists, nil
}

// ExistsByEmail e-postanın mevcut olup olmadığını kontrol eder
func (r *postgresUserRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("e-posta kontrolü yapılamadı: %w", err)
	}
	
	return exists, nil
} 