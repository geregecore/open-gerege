# Өөрчлөлтийн түүх (Changelog)

Энэ баримт бичигт төслийн бүх чухал өөрчлөлтүүдийг бүртгэнэ.

Формат нь [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) дээр суурилсан бөгөөд [Semantic Versioning](https://semver.org/spec/v2.0.0.html) дагадаг.

## [Хувилбаргүй] - Боловсруулалтын шат

### Нэмэгдсэн (Added)
- Анхны төслийн бүтэц үүсгэсэн
- Go backend Fiber framework дээр
- Next.js frontend App Router дээр
- PostgreSQL database schema
- Docker Compose тохиргоо

---

## [0.2.0] - 2025-01-22

### Нэмэгдсэн (Added)

#### Backend
- **Локал нэвтрэлт систем**
  - Email/password нэвтрэлт
  - Бүртгэлийн функц
  - Email баталгаажуулалт
  - Нууц үг сэргээх
  - MFA/TOTP дэмжлэг
  - Backup codes (random salt-тай)

- **Registration Service** (`internal/service/registration_service.go`)
  - `Register()` - Шинэ хэрэглэгч бүртгэх
  - `VerifyEmail()` - Email баталгаажуулах
  - `ResendVerificationEmail()` - Баталгаажуулалт дахин илгээх
  - `ForgotPassword()` - Нууц үг сэргээх хүсэлт
  - `ResetPassword()` - Нууц үг шинэчлэх

- **Registration Repository** (`internal/repository/registration_repo.go`)
  - Email verification token удирдлага
  - Password reset token удирдлага
  - User email verified статус

- **Registration Handler** (`internal/http/handlers/registration_handler.go`)
  - HTTP endpoints бүртгэл, баталгаажуулалтад

- **Шинэ API endpoints**
  - `POST /auth/local/register`
  - `POST /auth/local/verify-email`
  - `POST /auth/local/resend-verification`
  - `POST /auth/local/forgot-password`
  - `POST /auth/local/reset-password`

#### Frontend
- **Бүртгэлийн хуудас** (`src/app/register/page.tsx`)
- **RegisterForm** (`src/features/auth/components/RegisterForm.tsx`)
  - React Hook Form + Zod validation
  - Password strength indicator (zxcvbn)
  - Terms of service checkbox
  - WCAG 2.2 accessibility

- **PasswordStrengthIndicator** (`src/features/auth/components/PasswordStrengthIndicator.tsx`)
  - Нууц үгийн хүч харуулагч
  - 5 түвшний үнэлгээ

- **Email баталгаажуулах хуудас** (`src/app/verify-email/[token]/page.tsx`)
  - Token-ээр баталгаажуулалт
  - Success/Error төлөв

- **LogoutButton** (`src/features/auth/components/LogoutButton.tsx`)
  - Гарах товч component

- **Protected Routes Middleware** (`src/middleware.ts`)
  - Хамгаалагдсан routes
  - Session шалгалт
  - Auto redirect

- **Validation schemas** (`src/features/auth/schemas.ts`)
  - `loginSchema`
  - `registerSchema`
  - `forgotPasswordSchema`
  - `resetPasswordSchema`

- **Auth API updates** (`src/features/auth/api.ts`)
  - `register()`
  - `verifyEmail()`
  - `resendVerification()`
  - `forgotPassword()`
  - `resetPassword()`

#### Database
- **Шинэ migration бүтэц** (14 файл)
  - `001_extensions.sql` - PostgreSQL extensions
  - `002_core_tables.sql` - Systems, modules, permissions, roles
  - `003_user_tables.sql` - Users, citizens
  - `004_auth_tables.sql` - Credentials, MFA, sessions, tokens
  - `005_organization_tables.sql` - Organizations
  - `006_platform_tables.sql` - App icons, vehicles, devices
  - `007_content_tables.sql` - News, notifications, files, chat
  - `008_logging_tables.sql` - Audit, API logs, errors
  - `009_indexes.sql` - Performance indexes, full-text search
  - `010_seed_core.sql` - Systems, actions, modules
  - `011_seed_permissions.sql` - Permissions
  - `012_seed_roles.sql` - Roles, role permissions
  - `013_seed_organizations.sql` - Organizations, news categories
  - `014_seed_users.sql` - Admin users

#### Баримтжуулалт
- README.md Монгол хэл дээр
- backend/README.md шинэчилсэн
- frontend/README.md шинэчилсэн
- CONTRIBUTING.md нэмсэн
- CODE_OF_CONDUCT.md нэмсэн
- CHANGELOG.md нэмсэн
- SECURITY.md нэмсэн
- LICENSE файл нэмсэн

### Засварласан (Fixed)
- **Backup code salt аюулгүй байдал** - Hardcoded salt-ыг random salt болгон засав
  - `internal/service/auth_service.go` - `hashBackupCodeWithSalt()`, `generateBackupCodeSalt()`
  - `internal/domain/auth.go` - `UserMFABackupCode.Salt` field нэмсэн

- **LoginForm hardcoded credentials** - `admin@gerege.mn` / `admin123` устгасан
  - React Hook Form + Zod validation нэмсэн
  - ARIA labels нэмсэн

### Өөрчилсөн (Changed)
- Database schema `template_backend` болгож нэрлэсэн
- Database name `gerege_db` болгосон
- Migration файлуудыг дахин зохион байгуулсан

### Устгасан (Removed)
- Хуучин migration файлууд (001-013)
- LoginForm дахь hardcoded credentials

---

## [0.1.0] - 2025-01-21

### Нэмэгдсэн (Added)

#### Backend
- Go 1.25 Fiber v2 framework
- Clean Architecture бүтэц
- PostgreSQL GORM ORM
- Redis session cache
- SSO authentication
- RBAC permission систем
- OpenTelemetry tracing
- Prometheus metrics
- Swagger API docs
- Zap structured logging
- Health check endpoint
- Graceful shutdown

#### Frontend
- Next.js 16 App Router
- React 19
- TypeScript 5
- Tailwind CSS 4
- Zustand state management
- API client

#### Infrastructure
- Docker Compose тохиргоо
- GitHub Actions CI/CD
- Environment тохиргоо

### Аюулгүй байдал (Security)
- CORS тохиргоо
- CSP headers
- Rate limiting
- Input validation
- SQL injection хамгаалалт
- XSS хамгаалалт

---

## Хувилбарын тайлбар

- **Major (X.0.0)** - Таарахгүй өөрчлөлт (breaking changes)
- **Minor (0.X.0)** - Шинэ функц нэмэгдсэн (backwards compatible)
- **Patch (0.0.X)** - Алдаа засвар (backwards compatible)

## Зохиогчид

- Bayarsaikhan Otgonbayar
- Sengum
- Khuderchuluun
- Gankhulug

---

*Энэ changelog нь [Keep a Changelog](https://keepachangelog.com/) стандартыг дагадаг.*
