import { z } from 'zod';

// ============================================================
// LOGIN SCHEMA
// ============================================================

export const loginSchema = z.object({
    email: z
        .string()
        .min(1, 'Email хаяг оруулна уу')
        .email('Зөв email хаяг оруулна уу'),
    password: z
        .string()
        .min(1, 'Нууц үг оруулна уу'),
    rememberMe: z.boolean().optional(),
});

export type LoginFormData = z.infer<typeof loginSchema>;

// ============================================================
// REGISTRATION SCHEMA
// ============================================================

export const registerSchema = z.object({
    email: z
        .string()
        .min(1, 'Email хаяг оруулна уу')
        .email('Зөв email хаяг оруулна уу'),
    password: z
        .string()
        .min(8, 'Нууц үг хамгийн багадаа 8 тэмдэгт байх ёстой')
        .regex(/[A-Z]/, 'Нууц үг дор хаяж нэг том үсэг агуулсан байх ёстой')
        .regex(/[a-z]/, 'Нууц үг дор хаяж нэг жижиг үсэг агуулсан байх ёстой')
        .regex(/[0-9]/, 'Нууц үг дор хаяж нэг тоо агуулсан байх ёстой'),
    confirmPassword: z
        .string()
        .min(1, 'Нууц үг баталгаажуулна уу'),
    firstName: z
        .string()
        .min(1, 'Нэр оруулна уу')
        .max(150, 'Нэр 150 тэмдэгтээс хэтрэх ёсгүй'),
    lastName: z
        .string()
        .min(1, 'Овог оруулна уу')
        .max(150, 'Овог 150 тэмдэгтээс хэтрэх ёсгүй'),
    acceptTerms: z.literal(true, {
        errorMap: () => ({ message: 'Үйлчилгээний нөхцөлийг зөвшөөрнө үү' }),
    }),
}).refine((data) => data.password === data.confirmPassword, {
    message: 'Нууц үг таарахгүй байна',
    path: ['confirmPassword'],
});

export type RegisterFormData = z.infer<typeof registerSchema>;

// ============================================================
// FORGOT PASSWORD SCHEMA
// ============================================================

export const forgotPasswordSchema = z.object({
    email: z
        .string()
        .min(1, 'Email хаяг оруулна уу')
        .email('Зөв email хаяг оруулна уу'),
});

export type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

// ============================================================
// RESET PASSWORD SCHEMA
// ============================================================

export const resetPasswordSchema = z.object({
    password: z
        .string()
        .min(8, 'Нууц үг хамгийн багадаа 8 тэмдэгт байх ёстой')
        .regex(/[A-Z]/, 'Нууц үг дор хаяж нэг том үсэг агуулсан байх ёстой')
        .regex(/[a-z]/, 'Нууц үг дор хаяж нэг жижиг үсэг агуулсан байх ёстой')
        .regex(/[0-9]/, 'Нууц үг дор хаяж нэг тоо агуулсан байх ёстой'),
    confirmPassword: z
        .string()
        .min(1, 'Нууц үг баталгаажуулна уу'),
}).refine((data) => data.password === data.confirmPassword, {
    message: 'Нууц үг таарахгүй байна',
    path: ['confirmPassword'],
});

export type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

// ============================================================
// MFA VERIFICATION SCHEMA
// ============================================================

export const mfaVerificationSchema = z.object({
    code: z
        .string()
        .length(6, 'MFA код 6 оронтой байх ёстой')
        .regex(/^\d+$/, 'MFA код зөвхөн тоо агуулах ёстой'),
});

export type MfaVerificationFormData = z.infer<typeof mfaVerificationSchema>;

// ============================================================
// BACKUP CODE SCHEMA
// ============================================================

export const backupCodeSchema = z.object({
    code: z
        .string()
        .length(8, 'Backup код 8 тэмдэгт байх ёстой')
        .regex(/^[A-Z0-9]+$/, 'Backup код зөвхөн том үсэг, тоо агуулах ёстой'),
});

export type BackupCodeFormData = z.infer<typeof backupCodeSchema>;
