// Package domain provides implementation for domain
//
// File: auth.go
// Description: Authentication and security related domain models
// Author: Claude Code
// Created: 2026-01-11
/*
Package domain нь application-ийн бизнес entity-уудыг тодорхойлно.

Энэ файлд Authentication, MFA, Session, болон Security Audit entity-ууд
тодорхойлогдсон.

Database tables:
  - user_credentials: Local authentication password info
  - user_mfa_totp: TOTP MFA configuration
  - user_mfa_backup_codes: MFA recovery codes
  - sessions: User session tracking
  - login_history: Login attempt history
  - security_audit_trail: Security event logging
  - password_history: Password reuse prevention
*/
package domain

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// USER STATUS ENUM
// ============================================================

// UserStatus нь хэрэглэгчийн төлөвийг илэрхийлнэ.
type UserStatus string

const (
	// UserStatusActive - Идэвхтэй хэрэглэгч
	UserStatusActive UserStatus = "active"

	// UserStatusSuspended - Түр түдгэлзүүлсэн (админ үйлдэл)
	UserStatusSuspended UserStatus = "suspended"

	// UserStatusLocked - Түгжигдсэн (олон удаа буруу нууц үг оруулсан)
	UserStatusLocked UserStatus = "locked"

	// UserStatusPendingVerification - Баталгаажуулах хүлээгдэж буй
	UserStatusPendingVerification UserStatus = "pending_verification"

	// UserStatusDeactivated - Идэвхгүй болгосон
	UserStatusDeactivated UserStatus = "deactivated"
)

// IsValid checks if the status is a valid UserStatus
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusSuspended, UserStatusLocked,
		UserStatusPendingVerification, UserStatusDeactivated:
		return true
	}
	return false
}

// ============================================================
// USER CREDENTIAL ENTITY
// ============================================================

// UserCredential нь хэрэглэгчийн local authentication мэдээллийг хадгална.
// Table: user_credentials
//
// Password hash нь Argon2id алгоритм ашиглана.
// Account lockout: 5 удаа буруу нууц үг → 15 минут түгжээ.
type UserCredential struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"uniqueIndex;not null"`

	// PasswordHash нь Argon2id хэш
	PasswordHash string `json:"-" gorm:"not null"`

	// PasswordChangedAt нь нууц үг сүүлд солигдсон огноо
	PasswordChangedAt *time.Time `json:"password_changed_at"`

	// FailedLoginAttempts нь амжилтгүй нэвтрэлтийн тоо
	FailedLoginAttempts int `json:"failed_login_attempts" gorm:"default:0"`

	// LockedUntil нь account түгжигдсэн хугацаа
	LockedUntil *time.Time `json:"locked_until"`

	// MustChangePassword нь нэвтрэх үед нууц үг солих шаардлагатай эсэх
	MustChangePassword bool `json:"must_change_password" gorm:"default:false"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (UserCredential) TableName() string {
	return "user_credentials"
}

// IsLocked checks if the account is currently locked
func (uc *UserCredential) IsLocked() bool {
	if uc.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*uc.LockedUntil)
}

// ============================================================
// USER MFA TOTP ENTITY
// ============================================================

// UserMFATotp нь TOTP (Time-based One-Time Password) MFA тохиргоог хадгална.
// Table: user_mfa_totp
//
// Google Authenticator, Authy зэрэг app-тай нийцтэй.
// RFC 6238 стандарт.
type UserMFATotp struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"uniqueIndex;not null"`

	// SecretEncrypted нь шифрлэгдсэн TOTP secret
	// AES-256-GCM ашиглан шифрлэгдсэн
	SecretEncrypted string `json:"-" gorm:"not null"`

	// IsEnabled нь MFA идэвхтэй эсэх
	IsEnabled bool `json:"is_enabled" gorm:"default:false"`

	// VerifiedAt нь MFA баталгаажсан огноо
	VerifiedAt *time.Time `json:"verified_at"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (UserMFATotp) TableName() string {
	return "user_mfa_totp"
}

// ============================================================
// USER MFA BACKUP CODE ENTITY
// ============================================================

// UserMFABackupCode нь MFA сэргээх код хадгална.
// Table: user_mfa_backup_codes
//
// 10 ширхэг код үүсгэгдэх ба нэг удаа л ашиглаж болно.
type UserMFABackupCode struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"not null"`

	// CodeHash нь hash-лэгдсэн backup code
	CodeHash string `json:"-" gorm:"not null"`

	// Salt нь backup code hash-д ашиглагдсан random salt (base64 encoded)
	Salt string `json:"-" gorm:"type:varchar(64)"`

	// UsedAt нь код ашиглагдсан огноо (NULL бол ашиглаагүй)
	UsedAt *time.Time `json:"used_at"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (UserMFABackupCode) TableName() string {
	return "user_mfa_backup_codes"
}

// IsUsed checks if the backup code has been used
func (bc *UserMFABackupCode) IsUsed() bool {
	return bc.UsedAt != nil
}

// ============================================================
// SESSION ENTITY
// ============================================================

// Session нь хэрэглэгчийн session мэдээллийг хадгална.
// Table: sessions
//
// Redis нь primary storage, DB нь audit/backup.
type Session struct {
	// ID нь session ID (UUID)
	ID string `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"not null"`

	// IPAddress нь session үүсгэсэн IP
	IPAddress string `json:"ip_address"`

	// UserAgent нь browser/client мэдээлэл
	UserAgent string `json:"user_agent"`

	// ExpiresAt нь session дуусах хугацаа
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`

	// LastActivityAt нь сүүлийн үйлдлийн хугацаа
	LastActivityAt time.Time `json:"last_activity_at"`

	// RevokedAt нь session цуцлагдсан хугацаа
	RevokedAt *time.Time `json:"revoked_at"`

	// RevokedReason нь цуцлагдсан шалтгаан
	RevokedReason string `json:"revoked_reason"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRevoked checks if the session has been revoked
func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

// IsValid checks if the session is still valid (not expired and not revoked)
func (s *Session) IsValid() bool {
	return !s.IsExpired() && !s.IsRevoked()
}

// ============================================================
// LOGIN HISTORY ENTITY
// ============================================================

// LoginHistory нь нэвтрэлтийн түүхийг хадгална.
// Table: login_history
//
// Бүх нэвтрэлтийн оролдлого (амжилттай, амжилтгүй) бүртгэгдэнэ.
type LoginHistory struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key (nullable - амжилтгүй бол user олдоогүй байж болно)
	UserID *int `json:"user_id"`

	// Email нь нэвтрэхэд ашигласан email
	Email string `json:"email"`

	// IPAddress нь нэвтрэлтийн IP
	IPAddress string `json:"ip_address"`

	// UserAgent нь browser/client мэдээлэл
	UserAgent string `json:"user_agent"`

	// LoginMethod нь нэвтрэлтийн арга ('local', 'sso')
	LoginMethod string `json:"login_method" gorm:"not null"`

	// Success нь нэвтрэлт амжилттай эсэх
	Success bool `json:"success" gorm:"not null"`

	// FailureReason нь амжилтгүй болсон шалтгаан
	FailureReason string `json:"failure_reason"`

	// MFAUsed нь MFA ашигласан эсэх
	MFAUsed bool `json:"mfa_used" gorm:"default:false"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (LoginHistory) TableName() string {
	return "login_history"
}

// ============================================================
// SECURITY AUDIT TRAIL ENTITY
// ============================================================

// SecurityAuditAction нь аудит үйлдлийн төрлүүд
type SecurityAuditAction string

const (
	// Password actions
	AuditActionPasswordChange SecurityAuditAction = "password_change"
	AuditActionPasswordReset  SecurityAuditAction = "password_reset"

	// MFA actions
	AuditActionMFAEnable      SecurityAuditAction = "mfa_enable"
	AuditActionMFADisable     SecurityAuditAction = "mfa_disable"
	AuditActionMFABackupUsed  SecurityAuditAction = "mfa_backup_used"
	AuditActionMFABackupRegen SecurityAuditAction = "mfa_backup_regenerate"

	// Session actions
	AuditActionSessionCreate  SecurityAuditAction = "session_create"
	AuditActionSessionRevoke  SecurityAuditAction = "session_revoke"
	AuditActionSessionExpire  SecurityAuditAction = "session_expire"
	AuditActionLogoutAll      SecurityAuditAction = "logout_all"

	// Account actions
	AuditActionAccountLock    SecurityAuditAction = "account_lock"
	AuditActionAccountUnlock  SecurityAuditAction = "account_unlock"
	AuditActionStatusChange   SecurityAuditAction = "status_change"

	// Login actions
	AuditActionLoginSuccess SecurityAuditAction = "login_success"
	AuditActionLoginFailed  SecurityAuditAction = "login_failed"
)

// SecurityAuditTrail нь аюулгүй байдлын бүх үйлдлүүдийг бүртгэнэ.
// Table: security_audit_trail
type SecurityAuditTrail struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь үйлдэл хийсэн хэрэглэгчийн ID
	UserID *int `json:"user_id"`

	// Action нь үйлдлийн төрөл
	Action string `json:"action" gorm:"not null"`

	// TargetType нь зорилтот объектын төрөл
	TargetType string `json:"target_type"`

	// TargetID нь зорилтот объектын ID
	TargetID string `json:"target_id"`

	// OldValue нь өмнөх утга (JSON)
	OldValue string `json:"old_value" gorm:"type:jsonb"`

	// NewValue нь шинэ утга (JSON)
	NewValue string `json:"new_value" gorm:"type:jsonb"`

	// IPAddress нь үйлдлийн IP
	IPAddress string `json:"ip_address"`

	// UserAgent нь browser/client мэдээлэл
	UserAgent string `json:"user_agent"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (SecurityAuditTrail) TableName() string {
	return "security_audit_trail"
}

// ============================================================
// PASSWORD HISTORY ENTITY
// ============================================================

// PasswordHistory нь нууц үгийн түүхийг хадгална.
// Table: password_history
//
// Сүүлийн 5 нууц үгийг дахин ашиглахаас сэргийлнэ.
type PasswordHistory struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"not null"`

	// PasswordHash нь хуучин нууц үгийн hash
	PasswordHash string `json:"-" gorm:"not null"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (PasswordHistory) TableName() string {
	return "password_history"
}

// ============================================================
// USER STATUS EXTENSION
// ============================================================

// UserStatusInfo нь User model-д нэмэгдэх status талбарууд
// Энэ struct нь User model-д embed хийгдэнэ.
type UserStatusInfo struct {
	// Status нь хэрэглэгчийн төлөв
	Status string `json:"status" gorm:"default:active"`

	// StatusReason нь төлөв өөрчлөгдсөн шалтгаан
	StatusReason string `json:"status_reason"`

	// StatusChangedAt нь төлөв өөрчлөгдсөн огноо
	StatusChangedAt *time.Time `json:"status_changed_at"`

	// StatusChangedBy нь төлөв өөрчилсөн хэрэглэгчийн ID
	StatusChangedBy *int `json:"status_changed_by"`

	// LastLoginAt нь сүүлд нэвтэрсэн огноо
	LastLoginAt *time.Time `json:"last_login_at"`

	// LoginCount нь нийт нэвтэрсэн тоо
	LoginCount int `json:"login_count" gorm:"default:0"`
}

// ============================================================
// EMAIL VERIFICATION TOKEN ENTITY
// ============================================================

// EmailVerificationToken нь email баталгаажуулах токен хадгална.
// Table: email_verification_tokens
type EmailVerificationToken struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"not null"`

	// Token нь unique token string
	Token string `json:"-" gorm:"uniqueIndex;not null"`

	// ExpiresAt нь токен дуусах хугацаа
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`

	// UsedAt нь токен ашиглагдсан хугацаа
	UsedAt *time.Time `json:"used_at"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// IsExpired checks if the token has expired
func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *EmailVerificationToken) IsUsed() bool {
	return t.UsedAt != nil
}

// ============================================================
// PASSWORD RESET TOKEN ENTITY
// ============================================================

// PasswordResetToken нь нууц үг сэргээх токен хадгална.
// Table: password_reset_tokens
type PasswordResetToken struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"not null"`

	// Token нь unique token string
	Token string `json:"-" gorm:"uniqueIndex;not null"`

	// ExpiresAt нь токен дуусах хугацаа
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`

	// UsedAt нь токен ашиглагдсан хугацаа
	UsedAt *time.Time `json:"used_at"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired checks if the token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// ============================================================
// REFRESH TOKEN ENTITY
// ============================================================

// RefreshToken нь refresh token хадгална.
// Table: refresh_tokens
type RefreshToken struct {
	// ID нь primary key
	ID int `json:"id" gorm:"primaryKey"`

	// UserID нь users table руу foreign key
	UserID int `json:"user_id" gorm:"not null"`

	// TokenHash нь hash-лэгдсэн token
	TokenHash string `json:"-" gorm:"uniqueIndex;not null"`

	// SessionID нь session-тэй холбоотой
	SessionID string `json:"session_id" gorm:"not null"`

	// ExpiresAt нь токен дуусах хугацаа
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`

	// RevokedAt нь токен цуцлагдсан хугацаа
	RevokedAt *time.Time `json:"revoked_at"`

	// ExtraFields нь audit талбаруудыг агуулна
	ExtraFields

	// User нь холбогдсон хэрэглэгч
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
}

// TableName returns the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the token has expired
func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsRevoked checks if the token has been revoked
func (t *RefreshToken) IsRevoked() bool {
	return t.RevokedAt != nil
}

// ============================================================
// GORM HOOKS
// ============================================================

// BeforeCreate sets default values before creating
func (uc *UserCredential) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	uc.PasswordChangedAt = &now
	return nil
}
