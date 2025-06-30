#!/bin/bash

# Kütüphane Mikroservisleri Başlatma Scripti -   
# Bu script tüm    servisleri arka planda başlatır

echo "🚀 mikroservisleri başlatılıyor..."
echo "=================================================="

# Script dizinini al
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# PID dosyalarını saklamak için dizin oluştur
mkdir -p pids

# Önceki PID dosyalarını temizle
rm -f pids/*.pid

# Renkli çıktı için
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Fonksiyon:    servis başlat
start_service() {
    local service_name=$1
    local service_dir=$2
    local port=$3
    local   _type=$4
    
    echo -e "${BLUE}📦 $service_name başlatılıyor (Port: $port) - $  _type...${NC}"
    
    cd "$SCRIPT_DIR/$service_dir"
    
    #      servisleri için cmd/server dizinini kontrol et
    if [ -d "cmd/server" ]; then
        echo -e "${PURPLE}   🏗️      katmanları tespit edildi${NC}"
        cd "cmd/server"
        
        # Go modüllerini kontrol et (parent dizinde)
        if [ ! -f "../../go.mod" ]; then
            echo -e "${RED}❌ $service_dir dizininde go.mod bulunamadı!${NC}"
            return 1
        fi
        
        #      servisi başlat
        nohup go run main.go > "../../../logs/${service_name}.log" 2>&1 &
        local pid=$!
    else
        echo -e "${RED}❌      bekleniyor ama cmd/server bulunamadı!${NC}"
        return 1
    fi
    
    # PID'i dosyaya kaydet
    echo $pid > "$SCRIPT_DIR/pids/${service_name}.pid"
    
    # Servisin başladığını kontrol et
    sleep 3
    if kill -0 $pid > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $service_name başarıyla başladı (PID: $pid)${NC}"
        
        #      katmanlarını göster
        if [ "$  _type" = " " ]; then
            echo -e "${PURPLE}   🔗 Katmanlar: Domain → Use Cases → Interface Adapters → Frameworks${NC}"
        fi
        return 0
    else
        echo -e "${RED}❌ $service_name başlatılamadı!${NC}"
        return 1
    fi
}

# Fonksiyon: Port kontrolü
check_port() {
    local port=$1
    if netstat -an 2>/dev/null | grep ":$port " | grep LISTEN > /dev/null; then
        return 0
    else
        return 1
    fi
}

# Log dizini oluştur
mkdir -p logs

echo -e "${YELLOW}📋      Servisleri başlatma sırası:${NC}"
echo "1. 📖 Genre Service (Port: 3003) -     "
echo "2. ✍️  Author Service (Port: 3002) -     "  
echo "3. 📚 Book Service (Port: 3001) -     "
echo "4. 🔐 Auth Service (Port: 3005) -     "
echo "5. 🤖 Recommendation Service (Port: 3004) -     "
echo "6. 🌐 Gateway Service (Port: 3000) -     "
echo ""

# Servisleri sırayla başlat
failed_services=()

# 1. Genre Service ( )
if start_service "genre-service" "genre-service" "3003" " "; then
    sleep 3
    if ! check_port 3003; then
        echo -e "${YELLOW}⚠️  Genre Service portu henüz dinlemiyor, devam ediliyor...${NC}"
    fi
else
    failed_services+=("Genre Service")
fi

# 2. Author Service ( )
if start_service "author-service" "author-service" "3002" " "; then
    sleep 3
    if ! check_port 3002; then
        echo -e "${YELLOW}⚠️  Author Service portu henüz dinlemiyor, devam ediliyor...${NC}"
    fi
else
    failed_services+=("Author Service")
fi

# 3. Book Service ( )
if start_service "book-service" "book-service" "3001" " "; then
    sleep 3
    if ! check_port 3001; then
        echo -e "${YELLOW}⚠️  Book Service portu henüz dinlemiyor, devam ediliyor...${NC}"
    fi
else
    failed_services+=("Book Service")
fi

# 4. Auth Service ( )
if start_service "auth-service" "auth-service" "3005" " "; then
    sleep 3
    if ! check_port 3005; then
        echo -e "${YELLOW}⚠️  Auth Service portu henüz dinlemiyor, devam ediliyor...${NC}"
    fi
else
    failed_services+=("Auth Service")
fi

# 5. Recommendation Service ( )
if start_service "recommendation-service" "recommendation-service" "3004" ""; then 
    sleep 3
    if ! check_port 3004; then
        echo -e "${YELLOW}⚠️  Recommendation Service portu henüz dinlemiyor, devam ediliyor...${NC}"
    fi
else
    failed_services+=("Recommendation Service")
fi

# 6. Gateway Service ( ) - en son
if start_service "gateway-service" "gateway-service" "3000" " "; then
    sleep 3
    if ! check_port 3000; then
        echo -e "${YELLOW}⚠️  Gateway Service portu henüz dinlemiyor, devam ediliyor...${NC}"
    fi
else
    failed_services+=("Gateway Service")
fi

echo ""
echo "=================================================="

# Sonuçları göster
if [ ${#failed_services[@]} -eq 0 ]; then
    echo -e "${GREEN}🎉 Tüm    mikroservisleri başarıyla başlatıldı!${NC}"
    echo ""
    echo -e "${BLUE}📊 Mikroservis Durumu (    ):${NC}"
    echo -e "${PURPLE}• Genre Service:         http://localhost:3003/health   📖 [ ]${NC}"
    echo -e "${PURPLE}• Author Service:        http://localhost:3002/health   ✍️  [ ]${NC}"
    echo -e "${PURPLE}• Book Service:          http://localhost:3001/health   📚 [ ]${NC}"
    echo -e "${PURPLE}• Auth Service:          http://localhost:3005/health   🔐 [ ]${NC}"
    echo -e "${PURPLE}• Recommendation Service: http://localhost:3004/health   🤖 [ ]${NC}"
    echo -e "${PURPLE}• Gateway Service:       http://localhost:3000/health   🌐 [ ]${NC}"
    echo ""
    echo -e "${BLUE}🌐    Gateway Endpoints:${NC}"
    echo -e "${PURPLE}• Gateway API:           http://localhost:3000/api/${NC}"
    echo -e "${PURPLE}• Books (Enriched):      http://localhost:3000/api/books/enriched${NC}"
    echo -e "${PURPLE}• Authors (Detailed):    http://localhost:3000/api/authors${NC}"
    echo -e "${PURPLE}• Genres (Detailed):     http://localhost:3000/api/genres${NC}"
    echo -e "${PURPLE}• Recommendations:       http://localhost:3000/api/recommendations${NC}"
    echo ""
    echo -e "${BLUE}🔐 Auth Service Endpoints:${NC}"
    echo -e "${PURPLE}• User Register:         http://localhost:3005/auth/register${NC}"
    echo -e "${PURPLE}• User Login:            http://localhost:3005/auth/login${NC}"
    echo -e "${PURPLE}• User Profile:          http://localhost:3005/auth/profile${NC}"
    echo -e "${PURPLE}• Token Validation:      http://localhost:3005/auth/validate${NC}"
    echo -e "${PURPLE}• Change Password:       http://localhost:3005/auth/change-password${NC}"
    echo ""
    echo -e "${BLUE}🏗️      Katmanları:${NC}"
    echo -e "${PURPLE}   📋 Domain Layer:       Entities, Value Objects, Domain Services${NC}"
    echo -e "${PURPLE}   💼 Use Case Layer:     Business Logic, Application Services${NC}"
    echo -e "${PURPLE}   🔌 Interface Layer:    Controllers, Repositories, Gateways${NC}"
    echo -e "${PURPLE}   🖥️  Framework Layer:   Gin, PostgreSQL, HTTP Client${NC}"
    echo ""
    echo -e "${YELLOW}📝 Log dosyaları: ./logs/ dizininde${NC}"
    echo -e "${YELLOW}🆔 PID dosyaları: ./pids/ dizininde${NC}"
    echo ""
    echo -e "${BLUE}💡 Servisleri durdurmak için: ./stop-services.sh${NC}"
else
    echo -e "${RED}❌ Bazı servisler başlatılamadı:${NC}"
    for service in "${failed_services[@]}"; do
        echo -e "${RED}   • $service${NC}"
    done
    echo ""
    echo -e "${YELLOW}📝 Detaylar için log dosyalarını kontrol edin: ./logs/${NC}"
    echo -e "${BLUE}💡 Başarısız servisleri tekrar başlatmak için scripti yeniden çalıştırın${NC}"
fi

echo "" 