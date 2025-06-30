#!/bin/bash

# Kütüphane Mikroservisleri Durdurma Scripti -         
# Bu script          mikroservisleri güvenli şekilde durdurur

echo "🛑          mikroservisleri durduruluyor..."
echo "=================================================="

# Script dizinini al
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Renkli çıktı için
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# PID dosyaları dizini kontrolü
if [ ! -d "pids" ]; then
    echo -e "${YELLOW}⚠️  PID dosyaları bulunamadı. Servisler çalışmıyor olabilir.${NC}"
    echo -e "${BLUE}💡 Manuel olarak kontrol etmek için: ps aux | grep 'go run'${NC}"
    exit 0
fi

# Fonksiyon:    servis durdur
stop_service() {
    local service_name=$1
    local pid_file="pids/${service_name}.pid"
    local     _type=$2
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        echo -e "${BLUE}📦 $service_name durduruluyor (PID: $pid) - $    _type...${NC}"
        
        # PID'in hala çalışıp çalışmadığını kontrol et
        if kill -0 $pid > /dev/null 2>&1; then
            # Önce SIGTERM gönder (nazikçe durdur)
            kill $pid
            
            # 5 saniye bekle
            local count=0
            while [ $count -lt 5 ] && kill -0 $pid > /dev/null 2>&1; do
                sleep 1
                count=$((count + 1))
            done
            
            # Hala çalışıyorsa SIGKILL gönder
            if kill -0 $pid > /dev/null 2>&1; then
                echo -e "${YELLOW}⚠️  $service_name nazikçe durmadı, zorla kapatılıyor...${NC}"
                kill -9 $pid
                sleep 1
            fi
            
            # Son kontrol
            if kill -0 $pid > /dev/null 2>&1; then
                echo -e "${RED}❌ $service_name durdurulamadı!${NC}"
                return 1
            else
                echo -e "${GREEN}✅ $service_name başarıyla durduruldu${NC}"
                
                #          bilgisi göster
                if [ "$    _type" = "   " ]; then
                    echo -e "${PURPLE}   🏗️          katmanları güvenli şekilde kapatıldı${NC}"
                fi
                
                rm -f "$pid_file"
                return 0
            fi
        else
            echo -e "${YELLOW}⚠️  $service_name zaten çalışmıyor${NC}"
            rm -f "$pid_file"
            return 0
        fi
    else
        echo -e "${YELLOW}⚠️  $service_name için PID dosyası bulunamadı${NC}"
        return 0
    fi
}

echo -e "${YELLOW}📋 Durdurulacak servisler (        ):${NC}"
echo -e "${PURPLE}• 🌐 Gateway Service [   ]${NC}"
echo -e "${PURPLE}• 🤖 Recommendation Service [   ]${NC}"
echo -e "${PURPLE}• 🔐 Auth Service [   ]${NC}"
echo -e "${PURPLE}• 📚 Book Service [   ]${NC}"
echo -e "${PURPLE}• ✍️  Author Service [   ]${NC}"
echo -e "${PURPLE}• 📖 Genre Service [   ]${NC}"
echo ""

# Servisleri ters sırada durdur (Gateway en önce)
stopped_services=()
failed_services=()

# 1. Gateway Service (   ) - diğer servislere bağımlı olduğu için önce
if stop_service "gateway-service" "   "; then
    stopped_services+=("Gateway Service [   ]")
else
    failed_services+=("Gateway Service [   ]")
fi

# 2. Recommendation Service (   )
if stop_service "recommendation-service" "   "; then
    stopped_services+=("Recommendation Service [   ]")
else
    failed_services+=("Recommendation Service [   ]")
fi

# 3. Auth Service (   )
if stop_service "auth-service" "   "; then
    stopped_services+=("Auth Service [   ]")
else
    failed_services+=("Auth Service [   ]")
fi

# 4. Book Service (   )
if stop_service "book-service" "   "; then
    stopped_services+=("Book Service [   ]")
else
    failed_services+=("Book Service [   ]")
fi

# 5. Author Service (   )
if stop_service "author-service" "   "; then
    stopped_services+=("Author Service [   ]")
else
    failed_services+=("Author Service [   ]")
fi

# 6. Genre Service (   )
if stop_service "genre-service" "   "; then
    stopped_services+=("Genre Service [   ]")
else
    failed_services+=("Genre Service [   ]")
fi

echo ""
echo "=================================================="

# Manuel olarak çalışan Go processlerini kontrol et
echo -e "${BLUE}🔍 Manuel çalışan Go processlerini kontrol ediliyor...${NC}"
running_processes=$(ps aux | grep '[g]o run main.go' | wc -l)

if [ $running_processes -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Hala çalışan Go processler bulundu:${NC}"
    ps aux | grep '[g]o run main.go' | awk '{print "   PID: " $2 " - " $11 " " $12}'
    echo ""
    echo -e "${BLUE}💡 Manuel olarak durdurmak için: kill -9 <PID>${NC}"
    
    #          process'leri tanımla
    echo -e "${PURPLE}🏗️          process'leri cmd/server/main.go şeklinde çalışır${NC}"
else
    echo -e "${GREEN}✅ Çalışan Go process bulunamadı${NC}"
fi

echo ""

# Sonuçları göster
if [ ${#failed_services[@]} -eq 0 ]; then
    echo -e "${GREEN}🎉 Tüm          mikroservisleri başarıyla durduruldu!${NC}"
    
    echo -e "${BLUE}🏗️ Durdurulan          Katmanları:${NC}"
    echo -e "${PURPLE}   📋 Domain Layer:       Entity ve Value Object'ler temizlendi${NC}"
    echo -e "${PURPLE}   💼 Use Case Layer:     Business Logic bağlantıları kapatıldı${NC}"
    echo -e "${PURPLE}   🔌 Interface Layer:    HTTP Handler'lar ve Repository'ler kapatıldı${NC}"
    echo -e "${PURPLE}   🖥️  Framework Layer:   Gin, PostgreSQL bağlantıları güvenli şekilde kapatıldı${NC}"
    echo ""
    
    # PID klasörünü temizle
    if [ -d "pids" ] && [ -z "$(ls -A pids)" ]; then
        rmdir pids
        echo -e "${BLUE}📁 PID dizini temizlendi${NC}"
    fi
    
else
    echo -e "${RED}❌ Bazı servisler durdurulamadı:${NC}"
    for service in "${failed_services[@]}"; do
        echo -e "${RED}   • $service${NC}"
    done
    echo ""
    echo -e "${BLUE}💡 Başarısız servisleri manuel olarak durdurmak gerekebilir${NC}"
fi

if [ ${#stopped_services[@]} -gt 0 ]; then
    echo -e "${GREEN}✅ Başarıyla durdurulan servisler:${NC}"
    for service in "${stopped_services[@]}"; do
        echo -e "${GREEN}   • $service${NC}"
    done
fi

echo ""
echo -e "${BLUE}💡 Servisleri yeniden başlatmak için: ./start-services.sh${NC}"
echo "" 