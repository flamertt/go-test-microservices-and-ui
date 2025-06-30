# 📚 Kütüphane Yönetim Sistemi -    Mikroservis Mimarisi

Bu proje, **Go** ve **React** teknolojileri kullanılarak geliştirilmiş tam özellikli bir kütüphane yönetim sistemidir. **  ** prensipleri ile tasarlanmış **gerçek mikroservis mimarisi** kullanarak, her servis kendi sorumluluğunda veri sunar ve **servisler arası iletişim** kurarak zenginleştirilmiş veri sağlar.

## 🏗️    Sistem Mimarisi

### **   Diyagramı**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                             CLIENT LAYER                                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐               │
│  │   React Web UI  │  │   Mobile Apps   │  │  Other Apps  │               │
│  │  (Port: 5173)   │  │   (Future)      │  │   (Future)   │               │
│  │ • Vanilla CSS   │  │                 │  │              │               │
│  │ • Modern Design │  │                 │  │              │               │
│  │ • Glass Morphism│  │                 │  │              │               │
│  │ • Kompakt UI    │  │                 │  │              │               │
│  └─────────────────┘  └─────────────────┘  └──────────────┘               │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        API GATEWAY LAYER (  )              │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                    Gateway Service (Port: 3000)                        │ │
│  │  🏗️    Implementation:                                 │ │
│  │    📋 Domain: Proxy logic, Service discovery                           │ │
│  │    💼 Use Cases: Request routing, Health aggregation                   │ │
│  │    🔌 Interface: HTTP handlers, Service clients                        │ │
│  │    🖥️ Framework: Gin, HTTP client, CORS                               │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       MICROSERVICES LAYER                  │
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐│
│  │Book Service 📚 │  │Author Service ✍️│  │Genre Service 📖│  │Auth Service 🔐││
│  │(Port: 3001)    │◄─┤(Port: 3002)     │◄─┤(Port: 3003)     │  │(Port: 3005)  ││
│  │🏗️ Clean Arch   │  │🏗️ Clean Arch    │  │🏗️ Clean Arch    │  │🏗️ Clean Arch ││
│  │                │  │                 │  │                 │             │
│  │📋 Domain:      │  │📋 Domain:       │  │📋 Domain:       │  │📋 Domain:    ││
│  │• Book Entity   │  │• Author Entity  │  │• Genre Entity   │  │• User Entity ││
│  │• Search Params │  │• Book Info      │  │• Book Info      │  │• JWT Claims  ││
│  │• Domain Errors │  │• Domain Errors  │  │• Domain Errors  │  │• Auth Errors ││
│  │                │  │                 │  │                 │             │
│  │💼 Use Cases:   │  │💼 Use Cases:    │  │💼 Use Cases:    │  │💼 Use Cases: ││
│  │• GetBooks      │  │• GetAuthors     │  │• GetGenres      │  │• Register    ││
│  │• EnrichBooks   │  │• GetAuthorBooks │  │• GetGenreBooks  │  │• Login       ││
│  │                │  │                 │  │                 │  │• ValidateJWT ││
│  │                │  │                 │  │                 │             │
│  │🔌 Interface:   │  │🔌 Interface:    │  │🔌 Interface:    │  │🔌 Interface: ││
│  │• Repository    │  │• Repository     │  │• Repository     │  │• Repository  ││
│  │• HTTP Handler  │  │• HTTP Handler   │  │• HTTP Handler   │  │• HTTP Handler││
│  │• Author Service│  │• Book Service   │  │• Book Service   │  │• Middleware  ││
│  │                │  │                 │  │                 │             │
│  │🖥️ Framework:   │  │🖥️ Framework:    │  │🖥️ Framework:    │  │🖥️ Framework:││
│  │• Gin Router    │  │• Gin Router     │  │• Gin Router     │  │• Gin Router  ││
│  │• PostgreSQL    │  │• PostgreSQL     │  │• PostgreSQL     │  │• PostgreSQL  ││
│  │• HTTP Client   │  │• HTTP Client    │  │• HTTP Client    │  │• JWT/bcrypt  ││
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│           ▲                      ▲                   ▲                     │
│           │                      │                   │                     │
│           └──────────────────────┼───────────────────┘                     │
│                                  │                                         │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │               Recommendation Service (Port: 3004) [Legacy]             │ │
│  │               📄 Legacy Architecture (To be migrated)                  │ │
│  │         • Aggregates data from all    services         │ │
│  │         • Random recommendation algorithms                              │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DATA LAYER                                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  PostgreSQL DB  │  │  PostgreSQL DB  │  │  PostgreSQL DB  │             │
│  │    (Books)      │  │   (Authors)     │  │    (Genres)     │             │
│  │  67,000+ books  │  │   800+ authors  │  │   50+ genres    │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 🔧    Teknoloji Stack'i

### **Backend (   Mikroservisler)**
- **Programlama Dili:** Go 1.21+
- **Web Framework:** Gin Gonic (Framework Layer)
- **Veritabanı:** PostgreSQL (Framework Layer)
- **Architecture Pattern:**    (Uncle Bob)
- **Dependency Injection:** Constructor injection
- **Error Handling:** Domain-specific errors

### **   Prensipleri**
- **Dependency Inversion:** Inner layers don't depend on outer layers
- **Separation of Concerns:** Each layer has single responsibility
- **Testability:** Business logic independent of frameworks
- **Framework Independence:** Easy to swap Gin with other frameworks

### **Frontend (Modern React Uygulaması)**
- **Programlama Dili:** TypeScript
- **Framework:** React 18
- **Styling:** Vanilla CSS (Clean Design)
- **UI Features:** Glass morphism navbar, kompakt tasarım
- **Renk Temaları:** Mavi (Books), Yeşil (Authors), Turuncu (Genres), Mor (Recommendations)

## 📁    Proje Yapısı

```
library-microservices/
├── gateway-service/                 # 🌐 API Gateway [  ]
│   ├── cmd/server/main.go          # 🎯 Entry point + DI
│   ├── configs/config.go           # ⚙️ Configuration
│   ├── internal/
│   │   ├── handler/                # 🔌 Interface Layer
│   │   └── service/                # 💼 Use Case Layer
│   └── go.mod
│
├── book-service/                   # 📚 Book Service [  ]
│   ├── cmd/server/main.go          # 🎯 Entry point + DI
│   ├── configs/config.go           # ⚙️ Configuration
│   ├── internal/
│   │   ├── model/                  # 📋 Domain Layer
│   │   ├── repository/             # 🔌 Interface Layer
│   │   ├── service/                # 💼 Use Case Layer
│   │   └── handler/                # 🔌 Interface Layer
│   ├── data/books.go              # 🖥️ Framework Layer (Legacy)
│   ├── main.go.old                # Legacy backup
│   └── go.mod
│
├── author-service/                 # ✍️ Author Service [  ]
│   ├── cmd/server/main.go          # 🎯 Entry point + DI
│   ├── configs/config.go           # ⚙️ Configuration
│   ├── internal/
│   │   ├── model/                  # 📋 Domain Layer
│   │   ├── repository/             # 🔌 Interface Layer
│   │   ├── service/                # 💼 Use Case Layer
│   │   └── handler/                # 🔌 Interface Layer
│   ├── data/authors.go            # 🖥️ Framework Layer (Legacy)
│   ├── main.go.old                # Legacy backup
│   └── go.mod
│
├── genre-service/                  # 📖 Genre Service [  ]
│   ├── cmd/server/main.go          # 🎯 Entry point + DI
│   ├── configs/config.go           # ⚙️ Configuration
│   ├── internal/
│   │   ├── model/                  # 📋 Domain Layer
│   │   ├── repository/             # 🔌 Interface Layer
│   │   ├── service/                # 💼 Use Case Layer
│   │   └── handler/                # 🔌 Interface Layer
│   ├── data/genres.go             # 🖥️ Framework Layer (Legacy)
│   ├── main.go.old                # Legacy backup
│   └── go.mod
│
├── auth-service/                   # 🔐 Auth Service [  ]
│   ├── cmd/server/main.go          # 🎯 Entry point + DI
│   ├── configs/config.go           # ⚙️ Configuration
│   ├── internal/
│   │   ├── model/                  # 📋 Domain Layer
│   │   │   ├── user.go            # User entity & requests
│   │   │   └── errors.go          # Auth-specific errors
│   │   ├── repository/             # 🔌 Interface Layer
│   │   │   └── user_repository.go # PostgreSQL implementation
│   │   ├── service/                # 💼 Use Case Layer
│   │   │   └── auth_service.go    # Authentication logic
│   │   ├── handler/                # 🔌 Interface Layer
│   │   │   └── auth_handler.go    # HTTP endpoints
│   │   └── middleware/             # 🔌 Interface Layer
│   │       └── auth_middleware.go # JWT validation
│   ├── utils/                     # 🖥️ Framework Layer
│   │   ├── jwt.go                 # JWT token management
│   │   └── password.go            # Password hashing
│   ├── data/users.sql             # 🖥️ Database schema
│   └── go.mod
│
├── recommendation-service/          # 🤖 [Legacy - To be migrated]
│   ├── main.go                     # 📄 Legacy monolithic
│   └── go.mod
│
├── logs/                           # Servis logları
├── pids/                           # Process ID dosyaları
├── start-services.sh               # 🏗️    aware script
├── stop-services.sh                # 🏗️    aware script
└── README.md
```

## 🚀    Kurulum ve Çalıştırma

### **   Hızlı Başlangıç**
```bash
# Projeyi klonlayın
git clone <repo-url>
cd library-microservices

#    servisleri başlatın
./start-services.sh

# Enhanced output:
# 🚀    mikroservisleri başlatılıyor...
# 📖 Genre Service başlatılıyor (Port: 3003) -   ...
# 🏗️     katmanları tespit edildi
# ✅ Genre Service başarıyla başladı (PID: 12345)
# 🔗 Katmanlar: Domain → Use Cases → Interface Adapters → Frameworks

# Servisleri durdurmak için
./stop-services.sh
```

### **   Manuel Başlatma**
```bash
#    servisleri (cmd/server path)
cd auth-service/cmd/server && go run main.go       # Port 3005
cd genre-service/cmd/server && go run main.go      # Port 3003
cd author-service/cmd/server && go run main.go     # Port 3002
cd book-service/cmd/server && go run main.go       # Port 3001
cd gateway-service/cmd/server && go run main.go    # Port 3000

# Legacy servis
cd recommendation-service && go run main.go        # Port 3004
```

## 🔗    API Endpoint'leri

### **Gateway API (Port: 3000) -   **

#### **📚 Book Service (  )**
```bash
#    book endpoints
GET /api/books?page=1&page_size=50&search=kafka
GET /api/books/123                    # Enriched with author info
GET /api/books/enriched               # All books with author details
GET /api/books/author/Franz%20Kafka   # Books by author
GET /api/books/category/Literature    # Books by category
```

#### **✍️ Author Service (  )**
```bash
#    author endpoints
GET /api/authors?page=1&page_size=20&search=kafka
GET /api/authors/123                            # Author detail
GET /api/authors/detail/Franz%20Kafka           # Author + books
GET /api/authors/search?name=Franz              # Search authors
```

#### **📖 Genre Service (  )**
```bash
#    genre endpoints
GET /api/genres?page=1&page_size=20&search=literature
GET /api/genres/5                               # Genre detail
GET /api/genres/detail/Literature?page=1        # Genre + paginated books
GET /api/genres/search?name=science             # Search genres
```

#### **🤖 Recommendations (Legacy)**
```bash
# Legacy recommendation endpoints (to be migrated)
GET /api/recommendations?limit=15
GET /api/recommendations/category/Literature
GET /api/recommendations/author/Franz%20Kafka
```

### **🩺    Health Check**
```bash
curl http://localhost:3000/api/health

# Response:
{
  "gateway": "OK",
  "services": {
    "book-service": "OK [  ]",
    "author-service": "OK [  ]", 
    "genre-service": "OK [  ]",
    "recommendation-service": "OK [Legacy]"
  },
  "architecture": {
    "clean_architecture_services": 4,
    "legacy_services": 1,
    "migration_progress": "80%"
  }
}
```

## 🌐 Frontend - Modern Vanilla CSS Tasarım

### **🎨 Yeni Tasarım Özellikleri**

#### **Modern Navbar (Glass Morphism)**
- Şeffaf beyaz arka plan (95% opacity) + blur efekti
- Gradient brand icon (mor-mavi geçişli)
- 360° dönen animasyon hover'da
- Scroll-aware dinamik stil değişimi

#### **Kompakt Header & Search**
- Header boyutları %35-50 azaltıldı
- Title font-size: 3rem → 2.5rem
- Padding: 32px → 16px, margin: 40px → 24px
- Results header kaldırıldı (Books sayfası)

#### **Renk Temaları**
- **📚 Kitaplar:** Mavi tema (#1976d2, #2196f3)
- **✍️ Yazarlar:** Yeşil tema (#388e3c, #4caf50)
- **📖 Türler:** Turuncu tema (#e65100, #ff9800) - Yeni!
- **🤖 Öneriler:** Mor tema (#7b1fa2, #9c27b0)

#### **4 Buton Grid (Recommendations)**
- Grid layout: 4 buton yan yana
- Mobilde responsive: 2x2 → 1 sütun
- Butonlar: 🎲 Random, 📖 Kategori, ✍️ Yazar, 🔄 Yenile

## 💼    Use Case Örnekleri

### **Book Service - Enriched Book Detail**
```go
// Use Case Implementation
func (s *BookServiceImpl) GetEnrichedBookByID(id int) (*model.EnrichedBook, error) {
    // Domain validation
    if id < 1 {
        return nil, model.ErrInvalidBookID
    }
    
    // Repository call (Interface Adapter)
    book, err := s.bookRepo.GetBookByID(id)
    if err != nil {
        return nil, err
    }
    
    // External service call (Interface Adapter)
    authorInfo, err := s.authorService.GetAuthorInfo(book.Author)
    if err != nil {
        // Business rule: Continue without author info
        authorInfo = &model.AuthorInfo{Name: book.Author}
    }
    
    // Domain enrichment
    return book.ToEnriched(authorInfo), nil
}
```

### **Author Service - Author with Books**
```go
// Use Case: Get author with their books
func (s *AuthorServiceImpl) GetEnrichedAuthorByName(name string) (*model.EnrichedAuthor, error) {
    // Domain validation
    if name == "" {
        return nil, model.ErrInvalidAuthorName
    }
    
    // Repository call
    authors, err := s.authorRepo.GetAuthorByName(name)
    if err != nil {
        return nil, err
    }
    
    if len(authors) == 0 {
        return nil, model.ErrAuthorNotFound
    }
    
    author := authors[0]
    
    // Book Service integration
    books, err := s.bookService.GetBooksByAuthor(name)
    if err != nil {
        books = []model.BookInfo{} // Graceful fallback
    }
    
    return author.ToEnriched(books), nil
}
```

## 🔄    Mikroservis İletişimi

### **1. Book Detail + Author Info**
```
Client → Gateway → Book Service (Clean)
                ↓
Book Use Case → Author Service (Clean)
                ↓
Response: Enriched Book with Author Details
```

### **2. Author Detail + Books**
```
Client → Gateway → Author Service (Clean)
                ↓
Author Use Case → Book Service (Clean)
                ↓
Response: Author with Book List
```

### **3. Genre Detail + Paginated Books**
```
Client → Gateway → Genre Service (Clean)
                ↓
Genre Use Case → Book Service (Clean)
                ↓
Response: Genre with Paginated Books
```

## 📊    Migration Status

### **✅ Tamamlanan Servisler (  )**
- **Gateway Service** 🌐 - Proxy patterns, health aggregation
- **Book Service** 📚 - Domain entities, repository pattern, use cases
- **Author Service** ✍️ - Author domain, book integration
- **Genre Service** 📖 - Genre domain, book integration

### **🔄 Migration Bekleyen Servisler**
- **Recommendation Service** 🤖 - Legacy →    migration planned

### **Migration Progress: 80% (4/5 services)**

## 🧪 Test ve Monitoring

### **   Health Monitoring**
```bash
# Architecture-aware health check
curl http://localhost:3000/api/health

# Service-specific health
curl http://localhost:3001/health  # Book Service (Clean)
curl http://localhost:3002/health  # Author Service (Clean)
curl http://localhost:3003/health  # Genre Service (Clean)
curl http://localhost:3004/health  # Recommendation (Legacy)
```

### **Log Monitoring**
```bash
# All services
tail -f logs/*.log

#    services include layer info
# [CLEAN-ARCH] [BOOK-SERVICE] [USE-CASE] GetBooks called
# [CLEAN-ARCH] [AUTHOR-SERVICE] [REPOSITORY] Query executed
```

## 🚀 Gelecek Geliştirmeler

### **   Tamamlama**
- [ ] **Recommendation Service Migration** - Legacy →   
- [ ] **Domain Event System** - Event-driven communication
- [ ] **CQRS Pattern** - Command Query separation
- [ ] **Value Objects Enhancement** - Immutable domain objects

### **Infrastructure**
- [ ] **Docker   ** - Multi-stage builds per layer
- [ ] **Kubernetes** -    aware orchestration
- [ ] **API Gateway Enhancement** - Advanced routing
- [ ] **Monitoring** - Layer-specific metrics

### **Frontend**
- [ ] **Frontend   ** - React   
- [ ] **PWA Support** - Progressive Web App
- [ ] **Mobile App** - React Native

## 🤝 Katkıda Bulunma

### **   Guidelines**
1. **Domain Layer:** Sadece business entities ve kurallar
2. **Use Case Layer:** Framework'den bağımsız business logic
3. **Interface Adapters:** Dış dünya ile iletişim
4. **Framework Layer:** Gin, PostgreSQL implementasyonları

### **Contribution Process**
1. Fork yapın
2.    feature branch oluşturun
3. Layer separation'ı koruyun
4. Pull Request oluşturun

---

**🎯 Bu sistem, Uncle Bob'un    prensiplerini takip ederek geliştirilmiş, 4/5 mikroservis   'a geçirilmiş (%80 migration), modern frontend tasarım ile desteklenmiş profesyonel bir kütüphane yönetim sistemidir.** 