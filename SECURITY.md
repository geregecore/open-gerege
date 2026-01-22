# Аюулгүй байдлын бодлого

## Дэмжигдэж буй хувилбарууд

Одоогоор дараах хувилбарууд аюулгүй байдлын шинэчлэл авдаг:

| Хувилбар | Дэмжлэг |
| -------- | ------- |
| 0.2.x    | ✅ Дэмжигдэж байна |
| 0.1.x    | ❌ Дэмжлэг дууссан |
| < 0.1    | ❌ Дэмжлэг дууссан |

## Эмзэг байдал мэдэгдэх

Аюулгүй байдлын эмзэг байдал илрүүлсэн бол **нийтэд ил болгохгүйгээр** бидэнд мэдэгдэнэ үү.

### Мэдэгдэх арга

1. **Email**: security@gerege.mn
2. **Шифрлэлт**: PGP түлхүүр ашиглаж болно (доор үзнэ үү)

### Мэдэгдэлд оруулах мэдээлэл

- Эмзэг байдлын тодорхой тайлбар
- Дахин давтах алхмууд (Proof of Concept)
- Нөлөөлсөн хувилбар(ууд)
- Боломжит шийдэл (байвал)
- Таны холбоо барих мэдээлэл

### Хариу өгөх хугацаа

- **24 цагийн дотор** - Мэдэгдэл хүлээн авснаа баталгаажуулна
- **72 цагийн дотор** - Анхны үнэлгээ хийж, хариу өгнө
- **7 хоногийн дотор** - Засварын төлөвлөгөө мэдэгдэнэ
- **90 хоногийн дотор** - Засвар гарсны дараа нийтэд мэдэгдэнэ

## Аюулгүй байдлын арга хэмжээнүүд

### Backend

#### Нэвтрэлт таних
- Argon2id password hashing (санах ой: 64MB, давталт: 1, threads: 4)
- JWT token-д хугацаа тавих (access: 15 мин, refresh: 7 хоног)
- MFA/TOTP дэмжлэг
- Account lockout (5 удаа буруу оролдлого → 15 минут түгжих)
- Session удирдлага (нэг удаагийн logout, бүх session logout)

#### API аюулгүй байдал
- Rate limiting
- CORS тохиргоо
- CSP headers
- Input validation (бүх endpoint дээр)
- SQL injection хамгаалалт (GORM prepared statements)
- Request ID tracking

#### Өгөгдлийн хамгаалалт
- Мэдрэг өгөгдөл шифрлэх
- Soft delete pattern
- Audit logging
- Database-level encryption (PostgreSQL)

### Frontend

#### XSS хамгаалалт
- React-ийн автомат escaping
- `dangerouslySetInnerHTML` хэрэглэхгүй
- Content Security Policy

#### CSRF хамгаалалт
- SameSite cookie attribute
- CSRF token (form-д)

#### Мэдрэг өгөгдөл
- Local storage-д мэдрэг өгөгдөл хадгалахгүй
- HttpOnly, Secure cookies
- Нууц үгийг form-оос автоматаар устгах

### Infrastructure

#### Docker
- Non-root user ажиллуулах
- Read-only filesystem (боломжтой бол)
- Resource limits
- Security scanning (Trivy)

#### Database
- Хамгийн бага эрхийн зарчим
- Шифрлэгдсэн холболт (SSL)
- Тусдаа schema ашиглах

## Хамгаалалтын зөвлөмж

### Хөгжүүлэгчдэд

1. **Dependencies шинэчлэх** - `npm audit`, `go mod verify` тогтмол ажиллуулах
2. **Secrets удирдлага** - `.env` файлыг git-д оруулахгүй
3. **Code review** - Аюулгүй байдлын асуудлыг шалгах
4. **Testing** - Security тестүүд бичих

### Deployment хийхэд

1. **HTTPS** - Бүх холболтод TLS ашиглах
2. **Environment variables** - Production secrets тусдаа хадгалах
3. **Firewall** - Зөвхөн шаардлагатай портуудыг нээх
4. **Logging** - Бүх үйлдлийг бүртгэх
5. **Backup** - Тогтмол backup хийх

## Мэдэгдсэн эмзэг байдлууд

### CVE бүртгэл

Одоогоор мэдэгдсэн CVE байхгүй.

### Засарсан эмзэг байдлууд

| Огноо | Тайлбар | Засвар |
|-------|---------|--------|
| 2025-01-22 | Backup code hardcoded salt | v0.2.0-д random salt болгож засав |

## PGP түлхүүр

Шифрлэгдсэн мэдэгдэл илгээх бол дараах PGP түлхүүрийг ашиглана уу:

```
-----BEGIN PGP PUBLIC KEY BLOCK-----
(PGP түлхүүрийг энд оруулна)
-----END PGP PUBLIC KEY BLOCK-----
```

## Холбоо барих

- **Email**: security@gerege.mn
- **Response time**: 24-72 цаг

---

*Аюулгүй байдлын бодлого сүүлд шинэчлэгдсэн: 2025-01-22*
