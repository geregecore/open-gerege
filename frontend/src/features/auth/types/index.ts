// ============================================================
// LOGIN TYPES
// ============================================================

export interface LoginResponse {
    token?: string;
    access_token?: string;
    expires_at?: number;
    requires_mfa?: boolean;
    mfa_token?: string;
    user?: UserInfo;
}

export interface LocalLoginRequest {
    email: string;
    password: string;
}

export interface UserInfo {
    id: number;
    email: string;
    first_name: string;
    last_name: string;
    status: string;
}

// ============================================================
// REGISTRATION TYPES
// ============================================================

export interface RegisterRequest {
    email: string;
    password: string;
    confirmPassword: string;
    firstName: string;
    lastName: string;
    acceptTerms: boolean;
}

export interface RegisterResponse {
    userId: number;
    email: string;
    verificationSent: boolean;
    message: string;
}

// ============================================================
// EMAIL VERIFICATION TYPES
// ============================================================

export interface VerifyEmailRequest {
    token: string;
}

export interface VerifyEmailResponse {
    success: boolean;
    message: string;
}

export interface ResendVerificationRequest {
    email: string;
}

// ============================================================
// PASSWORD RESET TYPES
// ============================================================

export interface ForgotPasswordRequest {
    email: string;
}

export interface ResetPasswordRequest {
    token: string;
    password: string;
    confirmPassword: string;
}

export interface GenericResponse {
    success: boolean;
    message: string;
}

// ============================================================
// MFA TYPES
// ============================================================

export interface VerifyMFARequest {
    mfa_token: string;
    code: string;
}

export interface VerifyBackupCodeRequest {
    mfa_token: string;
    code: string;
}

// ============================================================
// SESSION TYPES
// ============================================================

export interface SessionInfo {
    session_id: string;
    ip_address: string;
    user_agent: string;
    created_at: string;
    last_active: string;
    is_current: boolean;
}

export interface SessionListResponse {
    sessions: SessionInfo[];
    total: number;
}
