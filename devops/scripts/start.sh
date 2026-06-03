#!/bin/bash

set -e
source devops/scripts/helper.sh

echo -e "\033c"
yago
echo -e "      🚀 ${GREEN}YAGO REST API - Quick Start${NC}"
echo ""

echo -e "🔳 Docker and Docker Compose - check..."
sleep 1
# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "\r\033[A❌ Docker and Docker Compose - ${RED}failed !!!${NC}\n"
    echo -e "${RED}❌ Docker is not installed${NC}"
    echo ""
    echo -e "Please install Docker first: ${BLUE}https://docs.docker.com/get-docker/${NC}"
    echo ""
    exit 0
fi
# Check Docker Compose
if ! docker compose version &> /dev/null 2>&1; then
    echo -e "\r\033[A❌ Docker and Docker Compose - ${RED}failed !!!${NC}\n"
    echo -e "${RED}❌ Docker Compose is not installed or running${NC}"
    echo ""
    echo -e "Please install Docker Compose: ${BLUE}https://docs.docker.com/compose/install/${NC} and run"
    echo ""
    exit 0
fi
echo -e "\r\033[A✅ Docker and Docker Compose - ${GREEN}installed${NC}"
# Create .env file if doesn't exist
if [ ! -f .env ]; then
    cp .env.example .env
    echo -e "${GREEN}📝 .env file created${NC}"
else
    echo -e "${GREEN}✅ .env file exists${NC}"
fi
echo "🔐 Checking JWT_SECRET..."
# Use make command to generate JWT_SECRET if missing
if ! grep -q "^JWT_SECRET=[^[:space:]]\{1,\}" .env 2>/dev/null; then
    if make generate-jwt-secret > /dev/null 2>&1; then
        echo -e "${GREEN}✅ JWT_SECRET generated and added to .env${NC}"
    else
        echo -e "${RED}❌ Failed to generate JWT_SECRET${NC}"
        echo -e "${YELLOW}Please run 'make generate-jwt-secret' manually to see the error${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ JWT_SECRET already configured${NC}"
fi

# Load environment variables from .env file
if [ -f .env ]; then
    echo -e "📄 Reading .env file..."
    set -a
    . .env
    set +a
    if [ ! -f ./api/.env ]; then
        cp .env ./api/.env
    fi
#    grep -E '^[A-Za-z_][A-Za-z0-9_]*=' .env | while IFS='=' read -r key value; do
#        echo "Loaded: $key=$value"
#    done
fi
if [ -f ./admin/.env.example ]; then
    if [ ! -f ./admin/.env ]; then
        cp ./admin/.env.example ./admin/.env
    fi
fi
echo -e "${GREEN}✅ .env file read${NC}"
SERVER_PORT=${SERVER_PORT:-9009}

# Generate SSL certificates only if they do not exist yet
CERT_DIR="./devops/nginx/ssl"
DOMAIN="${DOMAIN:-yago.loc}"
CERT_FILE="$CERT_DIR/$DOMAIN.crt"
KEY_FILE="$CERT_DIR/$DOMAIN.key"

if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    generate_ssl
else
    echo -e "${GREEN}✅ SSL certificates already exist${NC}"
fi

# Start containers
LOG_FILE="/tmp/yago-docker-compose.log"
trap cleanup_ui EXIT
# Stop existing containers if running
if docker compose ps | grep -q "Up"; then
    echo "Stopping existing containers..."
    docker compose down >"$LOG_FILE" 2>&1
fi

if run_with_animation docker compose up -d --build --wait; then
    print_success_screen
else
    print_error_screen
    exit 1
fi

echo ""
echo "💽 Running database migrations..."
# Run migrations with retry mechanism (database might need a moment after health check)
MAX_RETRIES=3
RETRY_DELAY=3

run_migrations() {
    local service=$1
    local command=$2
    local retries=0
    local success=false

    while [ $retries -lt $MAX_RETRIES ]; do
        if docker compose exec -T $service $command; then
            success=true
            break
        else
            retries=$((retries + 1))
            if [ $retries -lt $MAX_RETRIES ]; then
                echo -e "${YELLOW}⚠️  $service migration attempt $retries failed, retrying in ${RETRY_DELAY} seconds...${NC}"
                sleep $RETRY_DELAY
            fi
        fi
    done

    if [ "$success" = true ]; then
        echo -e "${GREEN}✅ $service migrations completed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to run $service migrations after $MAX_RETRIES attempts${NC}"
        exit 1
    fi
}

# Running migrations for the API (Go)
run_migrations "api" "go run cmd/migrate/main.go up"

# Running migrations for Admin (Laravel)
run_migrations "admin" "php artisan migrate --force"

# Running seeders for Admin (Laravel)
echo ""
echo "💾 Running database seeders..."
if docker compose exec -T admin php artisan db:seed --force; then
    echo -e "${GREEN}✅ Admin seeders completed successfully${NC}"
else
    echo -e "${RED}❌ Failed to run Admin seeders${NC}"
    exit 1
fi

# Running generate key for Admin (Laravel)
echo ""
echo "🗝️ Running generate key..."
if docker compose exec -T admin php artisan key:generate; then
    echo -e "${GREEN}✅ APP_KEY created successfully${NC}"
else
    echo -e "${RED}❌ Failed to run generate APP_KEY${NC}"
    exit 1
fi

echo ""
echo -e "${G1}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${G1}║${NC}         ${GREEN}🎉 Success! YAGO API is ready!${NC}           ${G1}║${NC}"
echo -e "${G1}╚══════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${G1}📍 YAGO API is running at:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}API Base:    https://api.$DOMAIN/api/v1${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Document UI: https://api.$DOMAIN/swagger/index.html${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Health:      https://api.$DOMAIN/health${NC}"
echo ""
echo -e "${G1}⚙️ YAGO ADMIN is running at:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Base URL:    https://admin.$DOMAIN${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Metrics UI:  https://admin.$DOMAIN/metrics${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Health:      https://admin.$DOMAIN/health${NC}"
echo ""
echo -e "${G1}📦 Docker Commands:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}View logs:   docker compose logs -f api${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Stop:        docker compose down${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Restart:     docker compose restart${NC}"
echo ""
echo -e "${G1}🛠️  Development Commands:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Run tests:   make test${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Run linter:  make lint${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Update docs: make swag${NC}"
echo ""
echo -e "${G1}💽  Database Commands:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Run migrations:     make migrate-up${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Rollback migration: make migrate-down${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Migration status:   make migrate-status${NC}"
echo ""
echo -e "${G1}👨‍💻 Admin Management:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Create admin:       make create-admin${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Promote user:       make promote-admin ID=<user_id>${NC}"
echo ""
echo -e "${G1}📄 Documentation:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}https://dch-zeon.github.io/yago-docs/${NC}"
echo ""
echo ""
echo -e "${G1}‍💻  Your credentials to admin:${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Login:     admin@mail.com${NC}"
echo -e "${RED}   ▷${NC} ${YELLOW}Password:  secret856${NC}"
echo ""
echo ""