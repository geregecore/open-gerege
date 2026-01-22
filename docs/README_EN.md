# Open-Gerege

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-16-black.svg)](https://nextjs.org/)

**Open-Gerege** is an enterprise-grade backend API template and web application framework built with Go and Next.js.

[**Монгол**](../README.md) | English

## Introduction

Open-Gerege is an open-source project built with Clean Architecture principles, containing all essential functions needed for modern web application development. Originally designed for Mongolian developers, it complies with international standards.

## Tech Stack

### Backend
| Technology | Version | Description |
|------------|---------|-------------|
| Go | 1.25 | Primary programming language |
| Fiber | v2 | High-performance web framework |
| PostgreSQL | 15+ | Primary database |
| Redis | 7+ | Session and cache management |
| GORM | v2 | ORM (Object-Relational Mapping) |
| OpenTelemetry | - | Observability, tracing |
| Swagger | - | API documentation |
| Zap | - | Structured logging |

### Frontend
| Technology | Version | Description |
|------------|---------|-------------|
| Next.js | 16 | React framework |
| React | 19 | UI library |
| TypeScript | 5 | Type-safe JavaScript |
| Tailwind CSS | 4 | Utility-first CSS |
| Zustand | - | State management |
| React Hook Form | - | Form management |
| Zod | - | Schema validation |

## Features

- **Authentication System**
  - SSO (Single Sign-On) integration
  - Local authentication (email/password)
  - MFA/TOTP support
  - Refresh token rotation

- **Access Control**
  - RBAC (Role-Based Access Control)
  - Fine-grained permission system
  - Organization-level access

- **User Management**
  - Registration, email verification
  - Password recovery
  - Profile management

- **Organization Management**
  - Multi-level organization structure
  - Employee management
  - Organization settings

- **Content Management**
  - News, notifications
  - File management
  - Chat system

- **Monitoring**
  - Health check endpoints
  - Prometheus metrics
  - API request logging
  - Audit logs

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+ (optional)
- Docker & Docker Compose (recommended)

### Quick Start with Docker

```bash
# Clone the repository
git clone https://github.com/geregecore/open-gerege.git
cd open-gerege

# Start with Docker Compose
docker-compose up -d
```

After startup:
- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- Swagger UI: http://localhost:8080/swagger/index.html

### Manual Installation

#### Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env file (database configuration)

# Install dependencies
go mod download

# Run database migrations
make migrate

# Start server
make run
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## Development

### Backend Commands

```bash
make run              # Run with live reload
make build            # Build binary
make test             # Run unit tests
make test-integration # Run integration tests
make swagger          # Generate Swagger docs
make lint             # Run linter
make audit            # Security audit
make mocks            # Generate mock files
```

### Frontend Commands

```bash
npm run dev           # Development server (Turbopack)
npm run build         # Production build
npm run start         # Production server
npm run lint          # Run ESLint
npm run type-check    # TypeScript check
```

## API Documentation

When the server is running, access Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Contributing

We welcome open-source contributions! Please read [CONTRIBUTING.md](../CONTRIBUTING.md) to get started.

## Security

If you discover a security vulnerability, please report it following the guidelines in [SECURITY.md](../SECURITY.md).

## License

This project is licensed under the MIT License. See [LICENSE](../LICENSE) for details.

## Authors

- **Bayarsaikhan Otgonbayar** - CTO, Gerege Core Team
- **Sengum** - Developer
- **Khuderchuluun** - Developer
- **Gankhulug** - Developer

## Contact

- **GitHub Issues**: [Report an issue](https://github.com/geregecore/open-gerege/issues)
- **Email**: info@gerege.mn

---

**Gerege Core Team** - Made for Mongolian developers 🇲🇳
