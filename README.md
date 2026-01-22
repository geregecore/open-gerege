# Open-Gerege

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-16-black.svg)](https://nextjs.org/)

**Open-Gerege** бол Go болон Next.js дээр суурилсан, байгууллагын түвшний backend API template болон вэб аппликейшн framework юм.

[English](./docs/README_EN.md) | **Монгол**

## Танилцуулга

Open-Gerege нь орчин үеийн вэб аппликейшн хөгжүүлэлтэд шаардлагатай бүх үндсэн функцуудыг агуулсан, Clean Architecture зарчмаар бүтээгдсэн нээлттэй эхийн төсөл юм. Монгол улсын хөгжүүлэгчдэд зориулан бүтээгдсэн боловч олон улсын стандартад нийцсэн.

## Технологийн стек

### Backend
| Технологи | Хувилбар | Тайлбар |
|-----------|----------|---------|
| Go | 1.25 | Үндсэн програмчлалын хэл |
| Fiber | v2 | Өндөр гүйцэтгэлтэй вэб framework |
| PostgreSQL | 15+ | Үндсэн өгөгдлийн сан |
| Redis | 7+ | Session болон cache удирдлага |
| GORM | v2 | ORM (Object-Relational Mapping) |
| OpenTelemetry | - | Observability, tracing |
| Swagger | - | API баримтжуулалт |
| Zap | - | Structured logging |

### Frontend
| Технологи | Хувилбар | Тайлбар |
|-----------|----------|---------|
| Next.js | 16 | React framework |
| React | 19 | UI library |
| TypeScript | 5 | Type-safe JavaScript |
| Tailwind CSS | 4 | Utility-first CSS |
| Zustand | - | State management |
| React Hook Form | - | Form удирдлага |
| Zod | - | Schema validation |

## Онцлог шинж чанарууд

- **Нэвтрэлт таних систем**
  - SSO (Single Sign-On) интеграци
  - Локал нэвтрэлт (email/password)
  - MFA/TOTP дэмжлэг
  - Refresh token rotation

- **Хандалтын удирдлага**
  - RBAC (Role-Based Access Control)
  - Нарийвчилсан permission систем
  - Байгууллагын түвшний хандалт

- **Хэрэглэгчийн удирдлага**
  - Бүртгэл, email баталгаажуулалт
  - Нууц үг сэргээх
  - Профайл удирдлага

- **Байгууллагын удирдлага**
  - Олон түвшний байгууллагын бүтэц
  - Ажилтны удирдлага
  - Байгууллагын тохиргоо

- **Контент удирдлага**
  - Мэдээ, мэдэгдэл
  - Файл удирдлага
  - Chat систем

- **Мониторинг**
  - Health check endpoints
  - Prometheus metrics
  - API request logging
  - Audit logs

## Төслийн бүтэц

```
open-gerege/
├── backend/                    # Go backend API
│   ├── cmd/server/            # Аппликейшн эхлэх цэг
│   ├── internal/
│   │   ├── app/               # Dependency injection
│   │   ├── auth/              # Authentication middleware
│   │   ├── config/            # Тохиргоо
│   │   ├── db/                # Database холболт
│   │   ├── domain/            # Domain entities
│   │   ├── http/
│   │   │   ├── dto/           # Data Transfer Objects
│   │   │   ├── handlers/      # HTTP handlers
│   │   │   └── router/        # Route тодорхойлолт
│   │   ├── middleware/        # HTTP middlewares
│   │   ├── repository/        # Data access layer
│   │   └── service/           # Business logic
│   ├── migrations/            # Database migrations
│   └── docs/                  # Swagger баримтжуулалт
│
├── frontend/                   # Next.js frontend
│   └── src/
│       ├── app/               # Next.js App Router
│       ├── components/        # UI components
│       ├── features/          # Feature modules
│       │   ├── auth/          # Authentication
│       │   └── ...
│       └── lib/               # Utilities
│
├── docs/                       # Төслийн баримтжуулалт
├── docker-compose.yml         # Docker тохиргоо
└── README.md                  # Энэ файл
```

## Эхлүүлэх

### Шаардлага

- Go 1.25+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+ (заавал биш)
- Docker & Docker Compose (санал болгох)

### Docker ашиглан эхлүүлэх (Хамгийн хялбар)

```bash
# Repository clone хийх
git clone https://github.com/geregecore/open-gerege.git
cd open-gerege

# Docker Compose ашиглан эхлүүлэх
docker-compose up -d
```

Үүний дараа:
- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- Swagger UI: http://localhost:8080/swagger/index.html

### Гараар суулгах

#### Backend

```bash
cd backend

# Environment файл хуулах
cp .env.example .env

# .env файлыг засах (database тохиргоо)

# Dependencies суулгах
go mod download

# Database migration ажиллуулах
make migrate

# Server эхлүүлэх
make run
```

#### Frontend

```bash
cd frontend

# Dependencies суулгах
npm install

# Development server эхлүүлэх
npm run dev
```

## Хөгжүүлэлт

### Backend командууд

```bash
make run              # Live reload-той ажиллуулах
make build            # Binary бүтээх
make test             # Unit тест ажиллуулах
make test-integration # Integration тест ажиллуулах
make swagger          # Swagger docs үүсгэх
make lint             # Linter ажиллуулах
make audit            # Security audit
make mocks            # Mock файлууд үүсгэх
```

### Frontend командууд

```bash
npm run dev           # Development server (Turbopack)
npm run build         # Production build
npm run start         # Production server
npm run lint          # ESLint ажиллуулах
npm run type-check    # TypeScript шалгах
```

## API баримтжуулалт

Server ажиллаж байх үед Swagger UI-г дараах хаягаар үзнэ:

```
http://localhost:8080/swagger/index.html
```

## Хувь нэмэр оруулах

Бид нээлттэй эхийн хамтын ажиллагааг дэмждэг! Хувь нэмэр оруулахын тулд [CONTRIBUTING.md](./CONTRIBUTING.md) файлыг уншина уу.

## Аюулгүй байдал

Аюулгүй байдлын асуудал илрүүлсэн бол [SECURITY.md](./SECURITY.md) файлд заасан журмын дагуу мэдэгдэнэ үү.

## Лиценз

Энэ төсөл MIT лицензийн дор түгээгдэж байна. Дэлгэрэнгүйг [LICENSE](./LICENSE) файлаас үзнэ үү.

## Зохиогчид

- **Bayarsaikhan Otgonbayar** - CTO, Gerege Core Team
- **Sengum** - Developer
- **Khuderchuluun** - Developer
- **Gankhulug** - Developer

## Холбоо барих

- **GitHub Issues**: [Асуудал мэдэгдэх](https://github.com/geregecore/open-gerege/issues)
- **Email**: info@gerege.mn

---

**Gerege Core Team** - Монголын хөгжүүлэгчдэд зориулав 🇲🇳
