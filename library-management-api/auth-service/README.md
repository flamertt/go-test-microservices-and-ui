# Auth Service

Bu mikroservis, library-microservices projesinin authentication ve authorization işlemlerini yönetir. JWT token tabanlı güvenli bir authentication sistemi sağlar.

## Özellikler

- ✅ Kullanıcı kaydı (Register)
- ✅ Kullanıcı girişi (Login) 
- ✅ JWT Token tabanlı authentication
- ✅ Token doğrulama ve yenileme
- ✅ Şifre güvenliği (bcrypt hashing)
- ✅ Güvenli middleware'lar
- ✅ PostgreSQL entegrasyonu
- ✅ CORS desteği

## Teknolojiler

- **Go 1.21**
- **Gin Framework** - Web framework
- **PostgreSQL** - Veritabanı
- **JWT** - Token tabanlı authentication
- **bcrypt** - Şifre hashing

## Kurulum

### 1. Dependency'leri yükle
```bash
cd auth-service
go mod tidy
```

### 2. Veritabanını hazırla
PostgreSQL'de `users` tablosunu oluşturun:
```sql
\i data/users.sql
```

### 3. Environment Variables
Aşağıdaki environment variable'ları ayarlayın:
```bash
export SERVER_PORT=3005
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_username
export DB_PASSWORD=your_password
export DB_NAME=your_database
export JWT_SECRET_KEY=your-super-secret-jwt-key
export JWT_TOKEN_DURATION=24h
```

### 4. Servisi başlat
```bash
go run cmd/server/main.go
```

## API Endpoints

### Public Endpoints (Authentication gerektirmez)

#### Kullanıcı Kaydı
```bash
POST /auth/register
Content-Type: application/json

{
    "username": "johndoe",
    "email": "john@example.com", 
    "password": "123456"
}
```

#### Kullanıcı Girişi
```bash
POST /auth/login
Content-Type: application/json

{
    "username": "johndoe",
    "password": "123456"
}
```

#### Token Yenileme
```bash
POST /auth/refresh
Content-Type: application/json

{
    "token": "your_jwt_token_here"
}
```

### Protected Endpoints (Authentication gerektirir)

Her istek için `Authorization` header'ında JWT token göndermelisiniz:
```
Authorization: Bearer your_jwt_token_here
```

#### Kullanıcı Profili
```bash
GET /auth/profile
Authorization: Bearer your_jwt_token_here
```

#### Şifre Değiştirme
```bash
POST /auth/change-password
Authorization: Bearer your_jwt_token_here
Content-Type: application/json

{
    "old_password": "123456",
    "new_password": "newpassword123"
}
```

#### Token Doğrulama
```bash
GET /auth/validate
Authorization: Bearer your_jwt_token_here
```

#### Kullanıcı Bilgisi (ID ile)
```bash
GET /auth/users/{id}
Authorization: Bearer your_jwt_token_here
```

## Diğer Mikroservislerle Entegrasyon

Bu auth-service'i diğer mikroservislerinizde authentication middleware olarak kullanabilirsiniz:

### 1. Token Doğrulama
```go
// Token doğrulama için auth-service'e istek
func validateToken(token string) (*UserClaims, error) {
    req, _ := http.NewRequest("GET", "http://localhost:3005/auth/validate", nil)
    req.Header.Set("Authorization", "Bearer " + token)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return nil, errors.New("invalid token")
    }
    
    var result struct {
        Valid    bool   `json:"valid"`
        UserID   uint   `json:"user_id"`
        Username string `json:"username"`
        Email    string `json:"email"`
    }
    
    json.NewDecoder(resp.Body).Decode(&result)
    
    return &UserClaims{
        UserID:   result.UserID,
        Username: result.Username,
        Email:    result.Email,
    }, nil
}
```

### 2. Middleware Örneği
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

## Örnek Kullanım

### 1. Kullanıcı Kaydı ve Girişi
```bash
# Kullanıcı kaydı
curl -X POST http://localhost:3005/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# Kullanıcı girişi  
curl -X POST http://localhost:3005/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

### 2. Token ile Korumalı Endpoint Erişimi
```bash
# Token al (login'den dönen token)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Profil bilgilerini al
curl -X GET http://localhost:3005/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

## Güvenlik

- ✅ Şifreler bcrypt ile hash'lenir
- ✅ JWT token'lar imzalanır
- ✅ Token süre sınırı (default: 24 saat)
- ✅ CORS koruması
- ✅ Input validation
- ✅ SQL injection koruması

## Health Check

Servisin çalışıp çalışmadığını kontrol etmek için:
```bash
curl http://localhost:3005/health
```

## Test

Servisi test etmek için:
```bash
go test ./...
```

## Port

Default port: **3005**

Diğer servisler:
- Author Service: 3002
- Book Service: 3001  
- Genre Service: 3003
- Gateway Service: 8080
- Recommendation Service: 3004

Bu auth-service'i diğer tüm mikroservislerinizle entegre ederek güvenli bir authentication sistemi kurabilirsiniz. 