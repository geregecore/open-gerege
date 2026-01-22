# Архитектурын баримтжуулалт

## Ерөнхий тойм

Open-Gerege нь Clean Architecture зарчмаар бүтээгдсэн, микросервис-бэлэн монолит архитектуртай систем юм.

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Client Layer                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │
│  │   Web App   │  │ Mobile App  │  │  Admin UI   │                  │
│  │  (Next.js)  │  │  (Flutter)  │  │  (Next.js)  │                  │
│  └─────────────┘  └─────────────┘  └─────────────┘                  │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼ HTTPS/REST
┌─────────────────────────────────────────────────────────────────────┐
│                          API Gateway                                 │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  Load Balancer → Rate Limiter → Auth → Router → Handlers   │    │
│  └─────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Backend Services                              │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐        │
│  │   Auth    │  │   User    │  │   Org     │  │  Content  │        │
│  │  Service  │  │  Service  │  │  Service  │  │  Service  │        │
│  └───────────┘  └───────────┘  └───────────┘  └───────────┘        │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Data Layer                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  PostgreSQL  │  │    Redis     │  │ File Storage │               │
│  │   (Primary)  │  │   (Cache)    │  │   (S3/Minio) │               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
└─────────────────────────────────────────────────────────────────────┘
```

## Backend архитектур

### Clean Architecture давхаргууд

```
┌─────────────────────────────────────────────┐
│             HTTP/Transport Layer            │
│  (Handlers, Middleware, DTOs, Router)       │
├─────────────────────────────────────────────┤
│              Service Layer                  │
│  (Business Logic, Use Cases)                │
├─────────────────────────────────────────────┤
│            Repository Layer                 │
│  (Data Access, Database Operations)         │
├─────────────────────────────────────────────┤
│              Domain Layer                   │
│  (Entities, Value Objects, Interfaces)      │
└─────────────────────────────────────────────┘
```

### Хавтасны бүтэц

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
│
├── internal/                    # Private packages
│   ├── app/
│   │   └── dependency.go        # DI container
│   │
│   ├── auth/
│   │   └── auth.go              # Auth middleware
│   │
│   ├── config/
│   │   └── config.go            # Configuration
│   │
│   ├── db/
│   │   └── db.go                # Database connection
│   │
│   ├── domain/                  # Domain entities
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── permission.go
│   │   ├── organization.go
│   │   └── auth.go
│   │
│   ├── http/
│   │   ├── dto/                 # Data Transfer Objects
│   │   │   ├── user_dto.go
│   │   │   └── auth_dto.go
│   │   │
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── user_handler.go
│   │   │   └── auth_handler.go
│   │   │
│   │   └── router/              # Route definitions
│   │       ├── router.go
│   │       └── auth_router.go
│   │
│   ├── middleware/              # HTTP middlewares
│   │   ├── logging.go
│   │   ├── recovery.go
│   │   └── cors.go
│   │
│   ├── repository/              # Data access
│   │   ├── user_repo.go
│   │   └── auth_repo.go
│   │
│   └── service/                 # Business logic
│       ├── user_service.go
│       └── auth_service.go
│
├── migrations/                  # SQL migrations
├── docs/                        # Swagger docs
└── docker/                      # Docker configs
```

### Dependency Injection

```go
// internal/app/dependency.go

type Container struct {
    // Repositories
    UserRepo    repository.UserRepository
    AuthRepo    repository.AuthRepository
    RegRepo     repository.RegistrationRepository

    // Services
    UserService    *service.UserService
    AuthService    *service.AuthService
    RegService     *service.RegistrationService

    // Handlers
    UserHandler    *handlers.UserHandler
    AuthHandler    *handlers.AuthHandler
    RegHandler     *handlers.RegistrationHandler
}

func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    authRepo := repository.NewAuthRepository(db)
    regRepo := repository.NewRegistrationRepository(db)

    // Initialize services
    userService := service.NewUserService(userRepo)
    authService := service.NewAuthService(authRepo, userRepo, cfg)
    regService := service.NewRegistrationService(authRepo, userRepo, regRepo, authService, cfg)

    // Initialize handlers
    userHandler := handlers.NewUserHandler(userService)
    authHandler := handlers.NewAuthHandler(authService)
    regHandler := handlers.NewRegistrationHandler(regService)

    return &Container{...}
}
```

## Frontend архитектур

### Feature-based бүтэц

```
frontend/src/
├── app/                         # Next.js App Router
│   ├── (auth)/                  # Auth route group
│   │   ├── login/
│   │   └── register/
│   ├── (dashboard)/             # Dashboard route group
│   │   ├── profile/
│   │   └── settings/
│   ├── layout.tsx
│   └── page.tsx
│
├── components/                  # Shared UI components
│   ├── ui/                      # Primitive UI elements
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   └── card.tsx
│   └── layout/                  # Layout components
│       ├── header.tsx
│       └── sidebar.tsx
│
├── features/                    # Feature modules
│   ├── auth/
│   │   ├── api.ts               # API calls
│   │   ├── components/          # Feature components
│   │   ├── hooks/               # Custom hooks
│   │   ├── schemas.ts           # Validation
│   │   ├── store.ts             # State
│   │   └── types/               # TypeScript types
│   │
│   ├── user/
│   └── organization/
│
├── lib/                         # Utilities
│   ├── api-client.ts
│   └── utils.ts
│
├── hooks/                       # Global hooks
├── stores/                      # Global state
└── middleware.ts                # Route middleware
```

### State Management (Zustand)

```typescript
// features/auth/store.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    setUser: (user: User | null) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            isAuthenticated: false,
            setUser: (user) => set({ user, isAuthenticated: !!user }),
            logout: () => set({ user: null, isAuthenticated: false }),
        }),
        { name: 'auth-storage' }
    )
);
```

## Database архитектур

### Schema diagram

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│     systems     │     │     modules     │     │   permissions   │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ id              │◄────│ system_id       │     │ id              │
│ name            │     │ id              │◄────│ module_id       │
│ code            │     │ parent_id       │     │ action_id       │────►
│ ...             │     │ name            │     │ name            │
└─────────────────┘     │ code            │     │ code            │
                        └─────────────────┘     └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│      roles      │     │role_permissions │     │    actions      │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ id              │◄────│ role_id         │     │ id              │
│ system_id       │     │ permission_id   │────►│ name            │
│ name            │     └─────────────────┘     │ code            │
│ code            │                             │ http_method     │
└─────────────────┘                             └─────────────────┘
        │
        ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   user_roles    │     │      users      │     │    citizens     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ user_id         │────►│ id              │────►│ id              │
│ role_id         │     │ citizen_id      │     │ register_number │
│ organization_id │     │ email           │     │ first_name      │
└─────────────────┘     │ first_name      │     │ last_name       │
                        │ status          │     │ ...             │
                        └─────────────────┘     └─────────────────┘
                                │
                                ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│user_credentials │     │  user_sessions  │     │  login_history  │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ user_id         │     │ user_id         │     │ user_id         │
│ credential_type │     │ session_id      │     │ login_type      │
│ password_hash   │     │ ip_address      │     │ status          │
│ ...             │     │ expires_at      │     │ ip_address      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

### Soft Delete Pattern

Бүх үндсэн хүснэгтүүд `deleted_date` талбартай:

```sql
-- Soft delete
UPDATE users SET deleted_date = NOW() WHERE id = 1;

-- Query (устгагдаагүй)
SELECT * FROM users WHERE deleted_date IS NULL;
```

### Audit Fields

Бүх хүснэгтүүд timestamp талбаруудтай:

```sql
created_date    TIMESTAMPTZ DEFAULT NOW()
updated_date    TIMESTAMPTZ DEFAULT NOW()
deleted_date    TIMESTAMPTZ              -- Soft delete
```

## Authentication Flow

### Локал нэвтрэлт

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │     │   API    │     │  Service │     │ Database │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │                │
     │  POST /login   │                │                │
     │───────────────►│                │                │
     │                │  Validate      │                │
     │                │───────────────►│                │
     │                │                │  Get User      │
     │                │                │───────────────►│
     │                │                │◄───────────────│
     │                │                │  Verify Pass   │
     │                │◄───────────────│                │
     │                │                │  Create Session│
     │                │                │───────────────►│
     │                │◄───────────────│                │
     │  Set Cookie    │                │                │
     │◄───────────────│                │                │
     │                │                │                │
```

### MFA Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │     │   API    │     │  Service │
└────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │
     │  POST /login   │                │
     │───────────────►│                │
     │                │  Check MFA     │
     │                │───────────────►│
     │  MFA Required  │                │
     │◄───────────────│                │
     │                │                │
     │  POST /mfa     │                │
     │───────────────►│                │
     │                │  Verify TOTP   │
     │                │───────────────►│
     │  Success       │                │
     │◄───────────────│                │
     │                │                │
```

## Caching Strategy

### Redis ашиглалт

```
┌─────────────────────────────────────────────────────────────┐
│                      Cache Layers                            │
├─────────────────────────────────────────────────────────────┤
│  Session Cache     │  TTL: 24h   │  user:{id}:session       │
│  User Cache        │  TTL: 1h    │  user:{id}:profile       │
│  Permission Cache  │  TTL: 15m   │  user:{id}:permissions   │
│  Rate Limit        │  TTL: 1m    │  ratelimit:{ip}          │
└─────────────────────────────────────────────────────────────┘
```

### Cache Invalidation

```go
// Service layer-т cache invalidation
func (s *UserService) UpdateUser(ctx context.Context, id int, data UpdateData) error {
    // Update database
    if err := s.repo.Update(ctx, id, data); err != nil {
        return err
    }

    // Invalidate cache
    s.cache.Delete(ctx, fmt.Sprintf("user:%d:profile", id))
    s.cache.Delete(ctx, fmt.Sprintf("user:%d:permissions", id))

    return nil
}
```

## Security Architecture

### Defense in Depth

```
┌─────────────────────────────────────────────────────────────┐
│ Layer 1: Network                                            │
│ • Firewall                                                  │
│ • DDoS Protection                                           │
│ • TLS 1.3                                                   │
├─────────────────────────────────────────────────────────────┤
│ Layer 2: Application                                        │
│ • Rate Limiting                                             │
│ • Input Validation                                          │
│ • CORS/CSP Headers                                          │
├─────────────────────────────────────────────────────────────┤
│ Layer 3: Authentication                                     │
│ • JWT Tokens                                                │
│ • MFA/TOTP                                                  │
│ • Session Management                                        │
├─────────────────────────────────────────────────────────────┤
│ Layer 4: Authorization                                      │
│ • RBAC                                                      │
│ • Resource-level permissions                                │
│ • Organization scoping                                      │
├─────────────────────────────────────────────────────────────┤
│ Layer 5: Data                                               │
│ • Encryption at rest                                        │
│ • Encryption in transit                                     │
│ • Audit logging                                             │
└─────────────────────────────────────────────────────────────┘
```

## Deployment Architecture

### Production Setup

```
                    ┌─────────────┐
                    │   Client    │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ CloudFlare  │
                    │    (CDN)    │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   Nginx     │
                    │ (Reverse    │
                    │   Proxy)    │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
    ┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
    │  Backend 1  │ │  Backend 2  │ │  Backend 3  │
    │   (Go)      │ │   (Go)      │ │   (Go)      │
    └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
           │               │               │
           └───────────────┼───────────────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
    ┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
    │ PostgreSQL  │ │    Redis    │ │   MinIO     │
    │  (Primary)  │ │  (Cluster)  │ │ (Storage)   │
    └─────────────┘ └─────────────┘ └─────────────┘
```

## Observability

### Logging

```go
// Structured logging with Zap
logger.Info("user created",
    zap.Int("user_id", user.ID),
    zap.String("email", user.Email),
    zap.String("request_id", ctx.Value("request_id").(string)),
)
```

### Metrics (Prometheus)

```
# HTTP metrics
http_requests_total{method="POST", path="/api/users", status="200"}
http_request_duration_seconds{method="POST", path="/api/users"}

# Database metrics
db_connections_open
db_connections_in_use
db_query_duration_seconds

# Business metrics
users_registered_total
login_attempts_total{status="success|failed"}
```

### Tracing (OpenTelemetry)

```go
// Distributed tracing
ctx, span := tracer.Start(ctx, "CreateUser")
defer span.End()

span.SetAttributes(
    attribute.Int("user_id", user.ID),
    attribute.String("email", user.Email),
)
```

---

*Архитектурын баримтжуулалт сүүлд шинэчлэгдсэн: 2025-01-22*
