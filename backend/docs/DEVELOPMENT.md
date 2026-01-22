# Development Guide

This guide helps developers set up and work with the backend codebase.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)
- Make

## Quick Start

```bash
# 1. Clone repository
git clone <repository-url>
cd backend-refactor-v25

# 2. Copy environment file
cp .env.example .env
# Edit .env with your settings

# 3. Start database (Docker)
docker-compose up -d postgres

# 4. Run migrations
make migrate

# 5. Start development server
make run
```

The server will be available at `http://localhost:8080`.

## Development Commands

```bash
# Run the application
make run

# Run with hot reload (requires air)
make dev

# Build binary
make build

# Run linter
make lint

# Format code
make fmt

# Generate Swagger docs
make swagger
```

## Testing

```bash
# Run all unit tests
make test-unit

# Run integration tests (requires Docker)
make test-integration

# Run all tests with coverage
make test-coverage

# View coverage report
open coverage.html
```

## Database

### Migrations

```bash
# Apply all pending migrations
make migrate

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create name=add_users_table
```

### Database Access

```bash
# Connect to database (Docker)
docker-compose exec postgres psql -U <user> -d <database>

# Reset database
make db-reset
```

## Code Organization

### Adding a New Feature

1. **Domain Entity** (`internal/domain/`)
   ```go
   // internal/domain/product.go
   type Product struct {
       ID        int       `gorm:"primaryKey"`
       Name      string    `gorm:"not null"`
       Price     float64   `gorm:"not null"`
       CreatedAt time.Time `gorm:"autoCreateTime"`
   }
   ```

2. **Repository Interface & Implementation** (`internal/repository/`)
   ```go
   // internal/repository/product_repo.go
   type ProductRepository interface {
       List(ctx context.Context, q dto.ProductListQuery) ([]domain.Product, int64, int, int, error)
       GetByID(ctx context.Context, id int) (domain.Product, error)
       Create(ctx context.Context, product domain.Product) error
       Update(ctx context.Context, id int, product domain.Product) error
       Delete(ctx context.Context, id int) error
   }
   ```

3. **Service Interface & Implementation** (`internal/service/`)
   ```go
   // internal/service/product_service.go
   type ProductService interface {
       List(ctx context.Context, q dto.ProductListQuery) ([]domain.Product, int64, int, int, error)
       Create(ctx context.Context, dto dto.ProductDto) error
       // ...
   }
   ```

4. **DTOs** (`internal/http/dto/`)
   ```go
   // internal/http/dto/product.go
   type ProductDto struct {
       Name  string  `json:"name" validate:"required,min=1,max=255"`
       Price float64 `json:"price" validate:"required,gt=0"`
   }
   ```

5. **Handler** (`internal/http/handlers/`)
   ```go
   // internal/http/handlers/product.go
   type ProductHandler struct {
       svc service.ProductService
   }

   func (h *ProductHandler) List(c *fiber.Ctx) error {
       // ...
   }
   ```

6. **Router** (`internal/http/router/`)
   ```go
   // internal/http/router/product_router.go
   func SetupProductRoutes(v1 fiber.Router, d *deps.AppDependencies) {
       h := handlers.NewProductHandler(d.ProductSvc)
       product := v1.Group("/product")
       product.Get("/", h.List)
       // ...
   }
   ```

7. **Wire Up** (`internal/http/wire.go`)
   - Register repository
   - Register service
   - Add router setup call

### Writing Tests

#### Unit Tests (Service Layer)

```go
// internal/service/product_service_test.go
func TestProductService_Create(t *testing.T) {
    tests := []struct {
        name      string
        input     dto.ProductDto
        mockSetup func(*mockProductRepository)
        wantErr   bool
    }{
        {
            name: "success - product created",
            input: dto.ProductDto{
                Name:  "Test Product",
                Price: 99.99,
            },
            mockSetup: func(m *mockProductRepository) {
                m.On("Create", mock.Anything, mock.AnythingOfType("domain.Product")).
                    Return(nil)
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockProductRepository{}
            tt.mockSetup(mockRepo)

            svc := NewProductService(mockRepo)
            err := svc.Create(context.Background(), tt.input)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

#### Integration Tests (Repository Layer)

```go
// tests/integration/product_repo_test.go
//go:build integration

package integration

func TestProductRepository_Create(t *testing.T) {
    ctx := context.Background()
    db := SetupTestDB(t)
    repo := repository.NewProductRepository(db)

    product := domain.Product{
        Name:  "Test Product",
        Price: 99.99,
    }

    err := repo.Create(ctx, product)
    assert.NoError(t, err)

    // Verify
    result, err := repo.GetByID(ctx, product.ID)
    assert.NoError(t, err)
    assert.Equal(t, "Test Product", result.Name)
}
```

### Generating Mocks

```bash
# Generate all mocks
make mocks

# Generate specific mock
mockery --dir=internal/repository --name=ProductRepository --output=tests/mocks
```

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing

### Naming Conventions

| Type | Convention | Example |
|------|-----------|---------|
| Package | lowercase | `repository` |
| Interface | CamelCase | `UserRepository` |
| Struct | CamelCase | `UserService` |
| Function | CamelCase | `GetByID` |
| Variable | camelCase | `userCount` |
| Constant | UPPER_SNAKE | `MAX_RETRIES` |
| JSON field | snake_case | `user_id` |

### Error Handling

```go
// Return errors, don't panic
func (s *UserService) GetByID(ctx context.Context, id int) (domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return domain.User{}, fmt.Errorf("get user by id: %w", err)
    }
    return user, nil
}
```

### Context Usage

Always pass context as the first parameter:

```go
func (r *userRepository) GetByID(ctx context.Context, id int) (domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return user, err
}
```

## Debugging

### Logging

```go
// Get logger from context
logger := logging.FromContext(ctx)

// Log with structured fields
logger.Info("user created",
    zap.Int("user_id", user.ID),
    zap.String("email", user.Email),
)
```

### Request Tracing

Each request has a unique ID in the `X-Request-ID` header. Use this for debugging:

```bash
# Find logs for a specific request
grep "request_id=abc123" logs/app.log
```

## API Documentation

### Swagger Annotations

```go
// @Summary Create a new user
// @Description Creates a new user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserDto true "User data"
// @Success 201 {object} resp.APIResponse[domain.User]
// @Failure 400 {object} resp.APIResponse[any]
// @Router /user [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
    // ...
}
```

### Regenerate Docs

```bash
make swagger
```

Access Swagger UI at: `http://localhost:8080/docs`

## Troubleshooting

### Common Issues

**1. Database connection failed**
```bash
# Check if database is running
docker-compose ps

# Check connection string in .env
echo $DB_HOST $DB_PORT
```

**2. Migration failed**
```bash
# Check migration status
make migrate-status

# Force a specific version
migrate -path migrations -database $DB_URL force <version>
```

**3. Tests failing**
```bash
# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestUserService_Create ./internal/service/...
```

**4. Lint errors**
```bash
# Run linter with details
golangci-lint run --verbose

# Fix auto-fixable issues
golangci-lint run --fix
```

## Performance Tips

1. **Use pagination** for list endpoints
2. **Add database indexes** for frequently queried columns
3. **Cache** frequently accessed data
4. **Use context timeouts** for external calls
5. **Profile** with pprof for performance bottlenecks

## Security Checklist

Before deploying, ensure:

- [ ] All endpoints have authentication
- [ ] Sensitive endpoints have rate limiting
- [ ] Input validation is in place
- [ ] SQL queries use parameterization
- [ ] Secrets are in environment variables
- [ ] HTTPS is enforced in production
- [ ] CORS is configured correctly
