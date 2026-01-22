# Open-Gerege Backend API

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7.svg)](https://gofiber.io/)

Clean Architecture зарчмаар бүтээгдсэн, өндөр гүйцэтгэлтэй Go backend API.

## Онцлог

- **Clean Architecture** - Domain-driven design, тодорхой хуваагдсан давхаргууд
- **Fiber v2** - Өндөр гүйцэтгэлтэй вэб framework
- **GORM** - PostgreSQL дэмжлэгтэй ORM
- **SSO интеграци** - Session caching-тэй Single Sign-On
- **Локал нэвтрэлт** - Email/password + MFA/TOTP
- **RBAC** - Role-Based Access Control
- **Observability** - OpenTelemetry tracing & Prometheus metrics
- **Integration Testing** - Testcontainers ашигласан тестүүд
- **CI/CD** - GitHub Actions pipeline
- **Swagger** - Автомат API баримтжуулалт
- **Structured Logging** - Zap logger, request ID дамжуулалт
- **Security Headers** - CSP, CORS, rate limiting
- **Graceful Shutdown** - Цэвэр resource cleanup

## Төслийн бүтэц

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Аппликейшн эхлэх цэг
├── internal/
│   ├── app/                     # Dependency injection container
│   ├── auth/                    # Authentication middleware
│   ├── config/                  # Тохиргооны бүтэц
│   ├── db/                      # Database холболт
│   ├── domain/                  # Domain models/entities
│   ├── http/
│   │   ├── dto/                 # Data transfer objects
│   │   ├── handlers/            # HTTP handlers
│   │   └── router/              # Route тодорхойлолт
│   ├── middleware/              # HTTP middlewares
│   ├── repository/              # Data access layer
│   └── service/                 # Business logic layer
├── migrations/                  # Database migrations
├── docs/                        # Swagger generated docs
└── docker/                      # Docker тохиргоо
```

## Түргэн эхлүүлэх

### Шаардлага

- Go 1.25+
- PostgreSQL 15+
- Redis 7+ (заавал биш)
- Make (заавал биш)

### Суулгалт

```bash
# Dependencies татах
go mod download

# Environment файл хуулах
cp .env.example .env

# .env файлыг засах

# Server ажиллуулах
go run cmd/server/main.go
```

### Make ашиглах

```bash
make run              # Server ажиллуулах
make build            # Binary бүтээх
make test             # Unit тест
make test-integration # Integration тест (Docker шаардана)
make test-all         # Бүх тест
make audit            # Security audit + linter
make mocks            # Mock үүсгэх (mockery)
make lint             # Linter
make swagger          # Swagger docs үүсгэх
make migrate          # Database migration
```

## Тохиргоо

`.env` файлын жишээ:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8000
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gerege_db
DB_SCHEMA=template_backend

# Redis (заавал биш)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Auth
AUTH_CACHE_TTL=1h
AUTH_CACHE_MAX=10000
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Local Auth
LOCAL_AUTH_ENABLED=true
PASSWORD_MIN_LENGTH=8
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m

# TLS (production-д)
TLS_CERT=
TLS_KEY=
```

## API Endpoints

### Нийтийн routes

| Method | Path | Тайлбар |
|--------|------|---------|
| GET | `/health` | Health check (DB статустай) |
| GET | `/swagger/*` | Swagger UI |

### Нэвтрэлт (Authentication)

| Method | Path | Тайлбар |
|--------|------|---------|
| GET | `/auth/login` | SSO redirect |
| GET | `/auth/callback` | OAuth2 callback |
| POST | `/auth/logout` | Гарах |
| GET | `/auth/verify` | Token шалгах |

### Локал нэвтрэлт

| Method | Path | Тайлбар |
|--------|------|---------|
| POST | `/auth/local/login` | Нэвтрэх |
| POST | `/auth/local/register` | Бүртгүүлэх |
| POST | `/auth/local/verify-email` | Email баталгаажуулах |
| POST | `/auth/local/resend-verification` | Баталгаажуулалт дахин илгээх |
| POST | `/auth/local/forgot-password` | Нууц үг сэргээх хүсэлт |
| POST | `/auth/local/reset-password` | Нууц үг шинэчлэх |
| POST | `/auth/local/refresh-token` | Token шинэчлэх |

### Хамгаалагдсан routes (нэвтрэлт шаардана)

| Resource | Endpoints | Тайлбар |
|----------|-----------|---------|
| User | `/user/*` | Хэрэглэгч |
| Role | `/role/*` | Роль |
| Permission | `/permission/*` | Зөвшөөрөл |
| Organization | `/organization/*` | Байгууллага |
| System | `/system/*` | Систем |
| Module | `/module/*` | Модуль |

## Health Check

`/health` endpoint нарийвчилсан статус буцаана:

```json
{
  "code": "OK",
  "data": {
    "status": "ok",
    "uptime": 3600,
    "timestamp": "2025-01-22T12:00:00Z",
    "database": {
      "status": "ok",
      "open_conns": 10,
      "in_use": 2,
      "idle": 8
    }
  }
}
```

## Хөгжүүлэлт

### Тест ажиллуулах

```bash
# Бүх тест
go test ./...

# Coverage-тэй
go test -v -race -coverprofile=coverage.out ./...

# Coverage тайлан харах
go tool cover -html=coverage.out
```

### Linting

```bash
# golangci-lint суулгах
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lint ажиллуулах
golangci-lint run
```

### Swagger Docs үүсгэх

```bash
# swag суулгах
go install github.com/swaggo/swag/cmd/swag@latest

# Docs үүсгэх
swag init -g cmd/server/main.go -o docs
```

## Database Migration

Migration файлууд `migrations/` хавтаст байрлана:

```
migrations/
├── 001_extensions.sql          # Extensions, functions
├── 002_core_tables.sql         # Systems, modules, permissions, roles
├── 003_user_tables.sql         # Users, citizens
├── 004_auth_tables.sql         # Credentials, sessions, tokens
├── 005_organization_tables.sql # Organizations
├── 006_platform_tables.sql     # App icons, vehicles, devices
├── 007_content_tables.sql      # News, notifications, files, chat
├── 008_logging_tables.sql      # Audit, API logs, errors
├── 009_indexes.sql             # Performance indexes
├── 010_seed_core.sql           # Systems, actions, modules seed
├── 011_seed_permissions.sql    # Permissions seed
├── 012_seed_roles.sql          # Roles seed
├── 013_seed_organizations.sql  # Organizations seed
└── 014_seed_users.sql          # Admin users seed
```

Migration ажиллуулах:

```bash
make migrate

# Эсвэл гараар
psql -U postgres -d gerege_db -f migrations/001_extensions.sql
# ... гэх мэт
```

## Docker

```bash
# Image бүтээх
docker build -t open-gerege-backend .

# Container ажиллуулах
docker run -p 8000:8000 --env-file .env open-gerege-backend
```

### Docker Compose

```bash
docker-compose up -d
```

## CI/CD

GitHub Actions workflows:

- **Lint** - golangci-lint шалгалт
- **Test** - PostgreSQL service-тэй unit тест
- **Build** - Binary compilation
- **Docker** - Image build & push
- **Security** - Gosec, Trivy scans

## Архитектур

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer                           │
│  ┌─────────┐  ┌─────────────┐  ┌──────────────────────┐    │
│  │ Router  │──│ Middleware  │──│      Handlers        │    │
│  └─────────┘  └─────────────┘  └──────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                          │
│  ┌──────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │ AuthService  │  │ UserService │  │ OrgService      │    │
│  └──────────────┘  └─────────────┘  └─────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Repository Layer                         │
│  ┌──────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │ AuthRepo     │  │ UserRepo    │  │ OrgRepo         │    │
│  └──────────────┘  └─────────────┘  └─────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│  ┌──────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │ User         │  │ Role        │  │ Organization    │    │
│  └──────────────┘  └─────────────┘  └─────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Layer                             │
│  ┌──────────────────────┐  ┌──────────────────────────┐    │
│  │     PostgreSQL       │  │         Redis            │    │
│  └──────────────────────┘  └──────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Лиценз

MIT License - [LICENSE](../LICENSE) файлаас дэлгэрэнгүй үзнэ үү.

## Зохиогчид

- Bayarsaikhan Otgonbayar, CTO - Gerege Core Team
- Sengum - Developer
- Khuderchuluun - Developer
- Gankhulug - Developer
