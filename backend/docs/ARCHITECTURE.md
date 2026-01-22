# Architecture Overview

This document describes the high-level architecture of the backend application.

## Layer Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Layer                                │
│  cmd/server/main.go → Router → Middleware → Handlers            │
├─────────────────────────────────────────────────────────────────┤
│                       Service Layer                              │
│  internal/service/*_service.go                                  │
│  (Business Logic, Validation, Orchestration)                    │
├─────────────────────────────────────────────────────────────────┤
│                     Repository Layer                             │
│  internal/repository/*_repo.go                                  │
│  (Data Access, GORM, Caching)                                   │
├─────────────────────────────────────────────────────────────────┤
│                       Domain Layer                               │
│  internal/domain/*.go                                           │
│  (Entities, Value Objects, Business Rules)                      │
└─────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── docs/
│   └── swagger.json             # OpenAPI specification
├── internal/
│   ├── domain/                  # Business entities
│   ├── http/
│   │   ├── dto/                 # Data Transfer Objects
│   │   ├── handlers/            # HTTP handlers
│   │   └── router/              # Route definitions
│   ├── middleware/              # HTTP middleware
│   ├── repository/              # Data access layer
│   └── service/                 # Business logic layer
├── migrations/                  # Database migrations
├── scripts/                     # Utility scripts
├── sql/                         # Raw SQL files
└── tests/
    ├── fixtures/                # Test seed data
    ├── integration/             # Integration tests
    ├── mocks/                   # Generated mocks
    └── testutils/               # Test utilities
```

## Dependency Flow

Dependencies flow inward only (Clean Architecture principle):

```
HTTP → Service → Repository → Domain
  │        │          │
  ▼        ▼          ▼
 DTO    Interface   GORM/DB
```

- **HTTP Layer** depends on Service interfaces
- **Service Layer** depends on Repository interfaces
- **Repository Layer** depends on Domain entities
- **Domain Layer** has no external dependencies

## Key Components

### 1. HTTP Layer

**Entry Point:** `cmd/server/main.go`
- Initializes configuration, database, and dependencies
- Sets up the Fiber HTTP server
- Applies middleware stack

**Router:** `internal/http/router/`
- Defines API routes and groups
- Applies authentication and authorization middleware
- Maps routes to handlers

**Handlers:** `internal/http/handlers/`
- Parse and validate HTTP requests
- Call service methods
- Format and return HTTP responses

**DTOs:** `internal/http/dto/`
- Request/response structures
- Validation rules (struct tags)
- JSON serialization configuration

### 2. Middleware Stack

Applied in order (see `internal/http/wire_security.go`):

1. **Recover** - Panic recovery
2. **Request ID** - Unique request identifier
3. **Helmet** - Security headers
4. **HSTS** - HTTP Strict Transport Security (production)
5. **HTTPS Redirect** - Force HTTPS (production)
6. **CORS** - Cross-Origin Resource Sharing
7. **CSRF** - Cross-Site Request Forgery protection
8. **Security Headers** - Additional security headers
9. **Body Size Limit** - Request body size limit (2MB)
10. **Rate Limiter** - Request rate limiting (100 req/min)
11. **Compression** - Response compression (gzip/brotli)
12. **Prometheus** - Metrics collection
13. **Request Context** - Context propagation
14. **Request Logger** - Access logging

### 3. Service Layer

**Location:** `internal/service/`

Services contain business logic and orchestrate operations:

```go
type UserService interface {
    List(ctx context.Context, q dto.UserListQuery) ([]domain.User, int64, int, int, error)
    Create(ctx context.Context, dto dto.UserDto) error
    Update(ctx context.Context, id int, dto dto.UserDto) error
    Delete(ctx context.Context, id int) error
}
```

Key responsibilities:
- Business rule validation
- Orchestrating repository calls
- Transaction management
- Domain event handling

### 4. Repository Layer

**Location:** `internal/repository/`

Repositories handle data persistence:

```go
type UserRepository interface {
    List(ctx context.Context, q dto.UserListQuery) ([]domain.User, int64, int, int, error)
    GetByID(ctx context.Context, id int) (domain.User, error)
    Create(ctx context.Context, user domain.User) error
    Update(ctx context.Context, id int, user domain.User) error
    Delete(ctx context.Context, id int) error
}
```

Key features:
- GORM for ORM operations
- Pagination support
- Query building
- Caching (where applicable)

### 5. Domain Layer

**Location:** `internal/domain/`

Domain entities represent business objects:

```go
type User struct {
    ID           int        `gorm:"primaryKey"`
    Email        string     `gorm:"unique;not null"`
    FirstName    string     `gorm:"not null"`
    LastName     string     `gorm:"not null"`
    RoleID       int        `gorm:"not null"`
    CreatedAt    time.Time  `gorm:"autoCreateTime"`
    UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}
```

## Authentication & Authorization

### Authentication

Uses SSO client for session-based authentication:

```go
// Middleware checks session validity
auth.RequireUser(cfg, ssoClient, sessionStore)
```

### Authorization

Permission-based access control:

```go
// Check specific permission
auth.RequirePermission(permCache, "admin.user.create")
```

Permission format: `{module}.{resource}.{action}`

## Database

- **ORM:** GORM v2
- **Database:** PostgreSQL 15+
- **Migrations:** Golang-migrate

### Connection Management

```go
// Database pool configuration
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)
```

## Observability

### Logging

- **Library:** Zap (structured logging)
- **Format:** JSON in production, console in development
- **Context:** Request ID, user ID, trace ID

### Metrics

- **Library:** Prometheus
- **Endpoint:** `/metrics`
- **Metrics:**
  - `http_requests_total` - Total HTTP requests
  - `http_request_duration_seconds` - Request latency
  - `db_query_duration_seconds` - Database query latency

### Health Checks

- **Endpoint:** `/health`
- **Checks:** Database connectivity, dependencies

## Security Features

| Feature | Implementation | Location |
|---------|---------------|----------|
| HTTPS | TLS termination | Load balancer / HSTS middleware |
| CSRF | Token-based | `middleware/csrf.go` |
| Rate Limiting | Per-user/IP | `middleware/limiter.go` |
| Input Validation | Struct tags | `http/dto/*.go` |
| SQL Injection | GORM (parameterized) | `repository/*.go` |
| XSS | Content-Type headers | `middleware/security.go` |

## API Design

### RESTful Conventions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/resource` | List resources |
| GET | `/resource/:id` | Get single resource |
| POST | `/resource` | Create resource |
| PUT | `/resource/:id` | Update resource |
| DELETE | `/resource/:id` | Delete resource |

### Response Format

**Success:**
```json
{
  "success": true,
  "code": "200",
  "msg": "success",
  "data": { ... },
  "pagination": {
    "page": 1,
    "size": 10,
    "total": 100
  }
}
```

**Error:**
```json
{
  "success": false,
  "code": "400",
  "msg": "validation_error",
  "data": null
}
```

## Testing Strategy

### Unit Tests
- Service layer with mocked repositories
- Location: `internal/service/*_test.go`

### Integration Tests
- Repository layer with real database (Testcontainers)
- Location: `tests/integration/*_test.go`

### Mocks
- Generated using mockery
- Location: `tests/mocks/`

## Configuration

Configuration is loaded from environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (dev/prod) | `development` |
| `APP_PORT` | HTTP port | `8000` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | - |
| `DB_USER` | Database user | - |
| `DB_PASSWORD` | Database password | - |
| `CORS_ORIGINS` | Allowed origins | `*` |

## Deployment

### Docker

```bash
# Build
docker build -t backend:latest .

# Run
docker run -p 8080:8080 --env-file .env backend:latest
```

### Health Check

```bash
curl http://localhost:8080/health
```
