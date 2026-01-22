import { apiClient } from '@/lib/api-client';
import {
    LocalLoginRequest,
    LoginResponse,
    RegisterRequest,
    RegisterResponse,
    VerifyEmailResponse,
    GenericResponse,
    ResetPasswordRequest,
    VerifyMFARequest,
    VerifyBackupCodeRequest,
} from './types/index';

const SSO_CONFIG = {
    origin: process.env.NEXT_PUBLIC_SSO_ORIGIN || 'https://sso.gerege.mn',
    clientId: process.env.NEXT_PUBLIC_SSO_CLIENT_ID || 'GRG-CLI-01KCGT4564YJ6WM15VNP3Y1BFG',
    redirectUri: process.env.NEXT_PUBLIC_REDIRECT_URI || 'http://localhost:3000',
};

export const getSSOLoginUrl = () => {
    const callbackUrl = `${SSO_CONFIG.redirectUri}/callback`;
    const encodedUri = encodeURIComponent(callbackUrl);
    return `${SSO_CONFIG.origin}/auth?client_id=${SSO_CONFIG.clientId}&redirect_uri=${encodedUri}`;
};

export const authApi = {
    // ============================================================
    // LOGIN / LOGOUT
    // ============================================================

    /**
     * Login with email and password
     */
    loginLocal: (data: LocalLoginRequest) =>
        apiClient.post<LoginResponse>('/auth/local/login', data),

    /**
     * Logout current session
     */
    logout: () =>
        apiClient.post<GenericResponse>('/auth/local/logout', {}),

    /**
     * Logout all sessions
     */
    logoutAll: () =>
        apiClient.post<GenericResponse>('/auth/local/logout-all', {}),

    // ============================================================
    // REGISTRATION
    // ============================================================

    /**
     * Register a new user
     */
    register: (data: RegisterRequest) =>
        apiClient.post<RegisterResponse>('/auth/local/register', {
            email: data.email,
            password: data.password,
            confirm_password: data.confirmPassword,
            first_name: data.firstName,
            last_name: data.lastName,
            accept_terms: data.acceptTerms,
        }),

    // ============================================================
    // EMAIL VERIFICATION
    // ============================================================

    /**
     * Verify email with token
     */
    verifyEmail: (token: string) =>
        apiClient.post<VerifyEmailResponse>('/auth/local/verify-email', { token }),

    /**
     * Resend verification email
     */
    resendVerification: (email: string) =>
        apiClient.post<GenericResponse>('/auth/local/resend-verification', { email }),

    // ============================================================
    // PASSWORD RESET
    // ============================================================

    /**
     * Request password reset email
     */
    forgotPassword: (email: string) =>
        apiClient.post<GenericResponse>('/auth/local/forgot-password', { email }),

    /**
     * Reset password with token
     */
    resetPassword: (data: ResetPasswordRequest) =>
        apiClient.post<GenericResponse>('/auth/local/reset-password', {
            token: data.token,
            password: data.password,
            confirm_password: data.confirmPassword,
        }),

    // ============================================================
    // MFA VERIFICATION
    // ============================================================

    /**
     * Verify MFA code (TOTP)
     */
    verifyMFA: (data: VerifyMFARequest) =>
        apiClient.post<LoginResponse>('/auth/local/verify-mfa', data),

    /**
     * Verify backup code
     */
    verifyBackupCode: (data: VerifyBackupCodeRequest) =>
        apiClient.post<LoginResponse>('/auth/local/verify-backup', data),

    // ============================================================
    // SESSION MANAGEMENT
    // ============================================================

    /**
     * Refresh current session
     */
    refreshToken: () =>
        apiClient.post<LoginResponse>('/auth/local/refresh', {}),
};
