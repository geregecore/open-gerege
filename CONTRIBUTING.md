# Хувь нэмэр оруулах заавар

Манай төсөлд сонирхож байгаад баярлалаа! Энэ баримт бичигт хувь нэмэр оруулах үйл явцыг тайлбарлав.

## Хүлээн авах зарчмууд

Бид бүх түвшний хөгжүүлэгчдийн хувь нэмрийг хүлээн авдаг. Та эхлэгч эсвэл туршлагатай хөгжүүлэгч байсан хамаагүй таны хувь нэмэр чухал.

## Хувь нэмэр оруулах арга замууд

### 1. Асуудал (Issue) мэдэгдэх

Алдаа, сайжруулах санал, шинэ функцын хүсэлт байвал GitHub Issues-ээр мэдэгдэнэ үү.

**Алдаа мэдэгдэхдээ:**
- Алдааг дахин давтах алхмууд
- Хүлээгдэж байсан үр дүн
- Бодит үр дүн
- Орчны мэдээлэл (OS, Go version, Node version гэх мэт)
- Боломжтой бол screenshot эсвэл log

**Шинэ функц хүсэхдээ:**
- Функцын тайлбар
- Хэрэглээний тохиолдол (use case)
- Боломжит шийдэл (заавал биш)

### 2. Код засвар (Pull Request)

#### Эхлэхээс өмнө

1. Ижил өөрчлөлт хийх гэж байгаа issue/PR байгаа эсэхийг шалгана уу
2. Том өөрчлөлтийн хувьд эхлээд issue нээж ярилцана уу
3. [Зан төлвийн дүрэм](./CODE_OF_CONDUCT.md)-тэй танилцана уу

#### Fork болон Clone

```bash
# Repository fork хийх (GitHub дээр)

# Өөрийн fork-оо clone хийх
git clone https://github.com/<таны-username>/open-gerege.git
cd open-gerege

# Upstream нэмэх
git remote add upstream https://github.com/geregecore/open-gerege.git
```

#### Branch үүсгэх

```bash
# Main branch-аас шинэ branch үүсгэх
git checkout main
git pull upstream main
git checkout -b feature/таны-feature-нэр
```

Branch нэрлэх дүрэм:
- `feature/xxx` - Шинэ функц
- `fix/xxx` - Алдаа засах
- `docs/xxx` - Баримтжуулалт
- `refactor/xxx` - Refactoring
- `test/xxx` - Тест нэмэх

#### Код бичих

**Backend (Go):**

```bash
cd backend

# Dependencies
go mod download

# Test ажиллуулах
go test ./...

# Lint
golangci-lint run

# Build
go build ./...
```

**Frontend (TypeScript):**

```bash
cd frontend

# Dependencies
npm install

# Lint
npm run lint

# Type check
npm run type-check

# Build
npm run build
```

#### Commit хийх

Бид [Conventional Commits](https://www.conventionalcommits.org/) стандартыг дагадаг:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Type-ууд:**
- `feat` - Шинэ функц
- `fix` - Алдаа засах
- `docs` - Баримтжуулалт
- `style` - Код форматлах (функц өөрчлөгдөөгүй)
- `refactor` - Refactoring
- `test` - Тест нэмэх/засах
- `chore` - Build, CI тохиргоо

**Жишээ:**

```bash
git commit -m "feat(auth): add email verification support"
git commit -m "fix(api): handle null response in user endpoint"
git commit -m "docs: update README with installation instructions"
```

#### Pull Request илгээх

1. Өөрчлөлтөө push хийх:
   ```bash
   git push origin feature/таны-feature-нэр
   ```

2. GitHub дээр Pull Request үүсгэх

3. PR template-ийг бөглөх:
   - Өөрчлөлтийн тайлбар
   - Холбогдох issue (байвал)
   - Тест хийсэн эсэх
   - Screenshot (UI өөрчлөлт байвал)

4. Code review хүлээх

### 3. Баримтжуулалт

Баримтжуулалтын сайжруулалт ч чухал хувь нэмэр юм:
- Алдаа засах
- Тайлбар нэмэх
- Жишээ нэмэх
- Орчуулга

## Код стандарт

### Go

- [Effective Go](https://golang.org/doc/effective_go) дагах
- `golangci-lint` ашиглах
- Test coverage 80%+
- Тодорхой нэрлэлт, comment

```go
// CreateUser creates a new user with the given parameters.
// It returns the created user and any error encountered.
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Implementation
}
```

### TypeScript

- ESLint + Prettier ашиглах
- Type-safe код
- React Hooks дүрэм дагах

```typescript
// Тодорхой type-тай
interface UserProps {
    id: number;
    name: string;
    email: string;
}

// Function component
export function UserCard({ id, name, email }: UserProps) {
    return (
        <div>
            <h2>{name}</h2>
            <p>{email}</p>
        </div>
    );
}
```

### SQL

- Table, column нэр lowercase, snake_case
- Index нэмэх (шаардлагатай газар)
- Migration файлд comment нэмэх

```sql
-- Create users table with audit fields
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(255) UNIQUE NOT NULL,
    created_date    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

## Review процесс

1. **Automated checks** - CI pipeline (lint, test, build)
2. **Code review** - Багийн гишүүдийн review
3. **Approval** - Хамгийн багадаа 1 approve
4. **Merge** - Maintainer merge хийнэ

## Тусламж авах

- GitHub Issues - Асуудал, санал
- GitHub Discussions - Ерөнхий ярилцлага
- Email: info@gerege.mn

## Тодруулга

Энэ төсөлд хувь нэмэр оруулснаар таны код MIT лицензийн дор түгээгдэх болохыг хүлээн зөвшөөрч байна.

---

Баярлалаа! 🙏
