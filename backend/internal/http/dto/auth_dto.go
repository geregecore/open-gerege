// Package dto provides implementation for dto
//
// File: auth_dto.go
// Description: DTOs for authentication, MFA, and session management
package dto

import "time"

// ============================================================
// LOGIN DTOs
// ============================================================

// LoginRequest нь local login хүсэлт
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse нь login хариу
type LoginResponse struct {
	RequiresMFA bool      `json:"requires_mfa,omitempty"`
	MFAToken    string    `json:"mfa_token,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresAt   int64     `json:"expires_at,omitempty"`
	User        *UserInfo `json:"user,omitempty"`
}

// UserInfo нь login хариунд буцаах хэрэглэгчийн мэдээлэл
type UserInfo struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
}

// ============================================================
// MFA DTOs
// ============================================================

// VerifyMFARequest нь MFA баталгаажуулах хүсэлт
type VerifyMFARequest struct {
	MFAToken string `json:"mfa_token" validate:"required"`
	Code     string `json:"code"      validate:"required,len=6"`
}

// VerifyBackupCodeRequest нь backup code баталгаажуулах хүсэлт
type VerifyBackupCodeRequest struct {
	MFAToken string `json:"mfa_token" validate:"required"`
	Code     string `json:"code"      validate:"required,len=8"`
}

// TOTPSetupResponse нь TOTP setup хариу
type TOTPSetupResponse struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
}

// ConfirmTOTPRequest нь TOTP баталгаажуулах хүсэлт
type ConfirmTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// DisableTOTPRequest нь TOTP идэвхгүй болгох хүсэлт
type DisableTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// MFAStatusResponse нь MFA төлөвийн хариу
type MFAStatusResponse struct {
	Enabled        bool `json:"enabled"`
	HasBackupCodes bool `json:"has_backup_codes"`
	BackupCodesLeft int `json:"backup_codes_left,omitempty"`
}

// BackupCodesResponse нь backup codes хариу
type BackupCodesResponse struct {
	Codes []string `json:"codes"`
}

// ============================================================
// PASSWORD DTOs
// ============================================================

// ChangePasswordRequest нь нууц үг солих хүсэлт
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
}

// SetPasswordRequest нь нууц үг тохируулах хүсэлт (admin)
type SetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

// ResetPasswordRequest нь нууц үг сэргээх хүсэлт
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ============================================================
// SESSION DTOs
// ============================================================

// SessionInfoResponse нь session мэдээллийн хариу
type SessionInfoResponse struct {
	SessionID  string    `json:"session_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
	IsCurrent  bool      `json:"is_current"`
}

// SessionListResponse нь session жагсаалтын хариу
type SessionListResponse struct {
	Sessions []SessionInfoResponse `json:"sessions"`
	Total    int                   `json:"total"`
}

// ============================================================
// USER STATUS DTOs
// ============================================================

// UpdateUserStatusRequest нь хэрэглэгчийн төлөв өөрчлөх хүсэлт
type UpdateUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active suspended locked deactivated"`
	Reason string `json:"reason" validate:"max=500"`
}

// ============================================================
// LOGIN HISTORY DTOs
// ============================================================

// LoginHistoryEntry нь login түүхийн оруулга
type LoginHistoryEntry struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	LoginMethod   string    `json:"login_method"`
	Success       bool      `json:"success"`
	FailureReason string    `json:"failure_reason,omitempty"`
	MFAUsed       bool      `json:"mfa_used"`
	CreatedAt     time.Time `json:"created_at"`
}

// LoginHistoryResponse нь login түүхийн хариу
type LoginHistoryResponse struct {
	Entries []LoginHistoryEntry `json:"entries"`
	Total   int                 `json:"total"`
}

// ============================================================
// SECURITY AUDIT DTOs
// ============================================================

// SecurityAuditEntry нь security audit оруулга
type SecurityAuditEntry struct {
	ID         int       `json:"id"`
	Action     string    `json:"action"`
	TargetType string    `json:"target_type,omitempty"`
	TargetID   string    `json:"target_id,omitempty"`
	OldValue   string    `json:"old_value,omitempty"`
	NewValue   string    `json:"new_value,omitempty"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// SecurityAuditResponse нь security audit хариу
type SecurityAuditResponse struct {
	Entries []SecurityAuditEntry `json:"entries"`
	Total   int                  `json:"total"`
}

// ============================================================
// REGISTRATION DTOs
// ============================================================

// RegisterRequest нь бүртгүүлэх хүсэлт
type RegisterRequest struct {
	Email           string `json:"email"            validate:"required,email"`
	Password        string `json:"password"         validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name"       validate:"required,min=1,max=150"`
	LastName        string `json:"last_name"        validate:"required,min=1,max=150"`
	AcceptTerms     bool   `json:"accept_terms"     validate:"required,eq=true"`
}

// RegisterResponse нь бүртгүүлэх хариу
type RegisterResponse struct {
	UserID           int    `json:"user_id"`
	Email            string `json:"email"`
	VerificationSent bool   `json:"verification_sent"`
	Message          string `json:"message"`
}

// VerifyEmailRequest нь email баталгаажуулах хүсэлт
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyEmailResponse нь email баталгаажуулах хариу
type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResendVerificationRequest нь verification email дахин илгээх хүсэлт
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ForgotPasswordRequest нь нууц үг мартсан хүсэлт
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordConfirmRequest нь нууц үг шинэчлэх хүсэлт
type ResetPasswordConfirmRequest struct {
	Token           string `json:"token"            validate:"required"`
	Password        string `json:"password"         validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// GenericResponse нь ерөнхий хариу (success/error message)
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
