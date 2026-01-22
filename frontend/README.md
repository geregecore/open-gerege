# Open-Gerege Frontend

[![Next.js](https://img.shields.io/badge/Next.js-16-black.svg)](https://nextjs.org/)
[![React](https://img.shields.io/badge/React-19-61DAFB.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6.svg)](https://www.typescriptlang.org/)

Next.js 16 болон React 19 дээр суурилсан орчин үеийн вэб аппликейшн.

## Технологи

| Технологи | Хувилбар | Тайлбар |
|-----------|----------|---------|
| Next.js | 16 | App Router, Server Components |
| React | 19 | UI library |
| TypeScript | 5 | Type-safe JavaScript |
| Tailwind CSS | 4 | Utility-first CSS framework |
| Zustand | - | Хөнгөн state management |
| React Hook Form | - | Form удирдлага |
| Zod | - | Schema validation |
| zxcvbn | - | Password strength шалгалт |

## Онцлог

- **App Router** - Next.js 16-ийн шинэ routing систем
- **Server Components** - Сервер талын rendering
- **TypeScript** - Бүрэн type safety
- **Responsive Design** - Mobile-first дизайн
- **Dark Mode** - Харанхуй/гэрэлтэй горим
- **WCAG 2.2** - Accessibility стандарт
- **Form Validation** - React Hook Form + Zod
- **Protected Routes** - Middleware-ээр хамгаалагдсан routes

## Төслийн бүтэц

```
frontend/
├── src/
│   ├── app/                     # Next.js App Router
│   │   ├── (auth)/              # Auth group routes
│   │   │   ├── login/           # Нэвтрэх хуудас
│   │   │   └── register/        # Бүртгүүлэх хуудас
│   │   ├── verify-email/        # Email баталгаажуулалт
│   │   ├── layout.tsx           # Root layout
│   │   └── page.tsx             # Нүүр хуудас
│   │
│   ├── components/              # Дахин ашиглагдах UI components
│   │   ├── ui/                  # Үндсэн UI elements
│   │   └── ...
│   │
│   ├── features/                # Feature modules
│   │   ├── auth/                # Authentication
│   │   │   ├── api.ts           # Auth API calls
│   │   │   ├── components/      # Auth components
│   │   │   │   ├── LoginForm.tsx
│   │   │   │   ├── RegisterForm.tsx
│   │   │   │   ├── LogoutButton.tsx
│   │   │   │   └── PasswordStrengthIndicator.tsx
│   │   │   ├── schemas.ts       # Zod schemas
│   │   │   └── types/           # TypeScript types
│   │   └── ...
│   │
│   ├── lib/                     # Utilities
│   │   ├── api-client.ts        # API client
│   │   └── utils.ts             # Helper functions
│   │
│   └── middleware.ts            # Next.js middleware
│
├── public/                      # Static files
├── tailwind.config.ts           # Tailwind тохиргоо
├── next.config.ts               # Next.js тохиргоо
└── package.json
```

## Эхлүүлэх

### Шаардлага

- Node.js 20+
- npm, yarn, pnpm, эсвэл bun

### Суулгалт

```bash
# Dependencies суулгах
npm install

# Development server эхлүүлэх
npm run dev
```

Хөтөч дээр http://localhost:2000 хаягаар нээнэ.

### Бусад командууд

```bash
npm run dev           # Development server (Turbopack)
npm run build         # Production build
npm run start         # Production server эхлүүлэх
npm run lint          # ESLint ажиллуулах
npm run type-check    # TypeScript шалгах
```

## Хуудсууд

| Path | Тайлбар | Хамгаалалт |
|------|---------|------------|
| `/` | Нүүр хуудас | Нийтийн |
| `/login` | Нэвтрэх | Нийтийн |
| `/register` | Бүртгүүлэх | Нийтийн |
| `/verify-email/[token]` | Email баталгаажуулах | Нийтийн |
| `/profile` | Профайл | Хамгаалагдсан |
| `/dashboard` | Хянах самбар | Хамгаалагдсан |
| `/settings` | Тохиргоо | Хамгаалагдсан |

## Authentication

### Нэвтрэх (Login)

```tsx
import { LoginForm } from '@/features/auth/components/LoginForm';

export default function LoginPage() {
  return <LoginForm />;
}
```

### Бүртгүүлэх (Register)

```tsx
import { RegisterForm } from '@/features/auth/components/RegisterForm';

export default function RegisterPage() {
  return <RegisterForm />;
}
```

### Гарах (Logout)

```tsx
import { LogoutButton } from '@/features/auth/components/LogoutButton';

export default function Header() {
  return <LogoutButton />;
}
```

## Form Validation

Zod schema ашиглан validation хийнэ:

```typescript
// features/auth/schemas.ts
import { z } from 'zod';

export const loginSchema = z.object({
  email: z.string()
    .min(1, 'Email оруулна уу')
    .email('Зөв email хаяг оруулна уу'),
  password: z.string()
    .min(1, 'Нууц үг оруулна уу'),
  rememberMe: z.boolean().optional(),
});

export const registerSchema = z.object({
  email: z.string()
    .min(1, 'Email оруулна уу')
    .email('Зөв email хаяг оруулна уу'),
  password: z.string()
    .min(8, 'Нууц үг хамгийн багадаа 8 тэмдэгт байх ёстой'),
  confirmPassword: z.string(),
  firstName: z.string()
    .min(1, 'Нэр оруулна уу')
    .max(150, 'Нэр хэт урт байна'),
  lastName: z.string()
    .min(1, 'Овог оруулна уу')
    .max(150, 'Овог хэт урт байна'),
  acceptTerms: z.literal(true, {
    errorMap: () => ({ message: 'Үйлчилгээний нөхцөл зөвшөөрөх шаардлагатай' }),
  }),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Нууц үг таарахгүй байна',
  path: ['confirmPassword'],
});
```

## API Client

```typescript
// lib/api-client.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export const apiClient = {
  async get<T>(endpoint: string): Promise<T> {
    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
      credentials: 'include',
    });
    if (!res.ok) throw new Error(res.statusText);
    return res.json();
  },

  async post<T>(endpoint: string, data: unknown): Promise<T> {
    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(res.statusText);
    return res.json();
  },
};
```

## Protected Routes Middleware

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const protectedRoutes = ['/profile', '/dashboard', '/settings'];
const authRoutes = ['/login', '/register'];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const sessionToken = request.cookies.get('session_token')?.value;

  // Хамгаалагдсан route - session шаардана
  if (protectedRoutes.some(route => pathname.startsWith(route))) {
    if (!sessionToken) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  // Auth routes - session байвал redirect
  if (authRoutes.includes(pathname)) {
    if (sessionToken) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    }
  }

  return NextResponse.next();
}
```

## Environment Variables

`.env.local` файл үүсгэх:

```env
# API
NEXT_PUBLIC_API_URL=http://localhost:8000

# App
NEXT_PUBLIC_APP_NAME=Open-Gerege
NEXT_PUBLIC_APP_URL=http://localhost:2000
```

## Styling

Tailwind CSS 4 ашиглана. `tailwind.config.ts` файлд тохиргоо байна:

```typescript
import type { Config } from 'tailwindcss';

const config: Config = {
  darkMode: 'class',
  content: [
    './src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        // ...
      },
    },
  },
};

export default config;
```

## Build & Deploy

### Production Build

```bash
npm run build
npm run start
```

### Docker

```bash
docker build -t open-gerege-frontend .
docker run -p 2000:2000 open-gerege-frontend
```

### Vercel

Vercel дээр deploy хийхэд хамгийн хялбар:

1. GitHub repo-г Vercel-д холбох
2. Environment variables тохируулах
3. Deploy товч дарах

## Лиценз

MIT License

## Зохиогчид

- Gerege Core Team
