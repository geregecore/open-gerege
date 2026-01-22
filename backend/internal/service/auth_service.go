// Package service provides implementation for service
//
// File: auth_service.go
// Description: Authentication service for local auth, MFA, and session management
package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"templatev25/internal/config"
	"templatev25/internal/domain"
	"templatev25/internal/repository"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

// Error definitions
var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrAccountLocked       = errors.New("account is locked")
	ErrAccountNotActive    = errors.New("account is not active")
	ErrMFARequired         = errors.New("MFA verification required")
	ErrInvalidMFACode      = errors.New("invalid MFA code")
	ErrMFANotEnabled       = errors.New("MFA is not enabled")
	ErrMFAAlreadyEnabled   = errors.New("MFA is already enabled")
	ErrInvalidSession      = errors.New("invalid or expired session")
	ErrPasswordTooWeak     = errors.New("password does not meet requirements")
	ErrPasswordReused      = errors.New("password was recently used")
	ErrUserNotFound        = errors.New("user not found")
	ErrCredentialsNotFound = errors.New("credentials not found")
)

// Argon2id parameters (OWASP recommended)
const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64MB
	argon2Threads = 4
	argon2KeyLen  = 32
	argon2SaltLen = 16
)

// AuthService handles authentication, MFA, and session management
type AuthService struct {
	repo         repository.AuthRepository
	sessionStore SessionStore
	cfg          *config.LocalAuthConfig
	logger       *zap.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(
	repo repository.AuthRepository,
	sessionStore SessionStore,
	cfg *config.LocalAuthConfig,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		repo:         repo,
		sessionStore: sessionStore,
		cfg:          cfg,
		logger:       logger,
	}
}

// ============================================================
// LOGIN
// ============================================================

// LoginRequest contains login parameters
type LoginRequest struct {
	Email     string
	Password  string
	IPAddress string
	UserAgent string
}

// LoginResponse contains login result
type LoginResponse struct {
	RequiresMFA bool
	MFAToken    string
	Session     *SessionData
	User        *domain.User
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logFailedLogin(ctx, nil, req.Email, req.IPAddress, req.UserAgent, "user not found")
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check user status
	if user.Status != string(domain.UserStatusActive) {
		s.logFailedLogin(ctx, &user.Id, req.Email, req.IPAddress, req.UserAgent, "account not active")
		return nil, ErrAccountNotActive
	}

	// Get credentials
	cred, err := s.repo.GetCredentialByUserID(ctx, user.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logFailedLogin(ctx, &user.Id, req.Email, req.IPAddress, req.UserAgent, "no credentials")
			return nil, ErrCredentialsNotFound
		}
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	// Check if account is locked
	if cred.IsLocked() {
		s.logFailedLogin(ctx, &user.Id, req.Email, req.IPAddress, req.UserAgent, "account locked")
		return nil, ErrAccountLocked
	}

	// Verify password
	if !s.verifyPassword(req.Password, cred.PasswordHash) {
		// Increment failed attempts
		s.repo.IncrementFailedAttempts(ctx, user.Id)

		// Check if should lock
		newCred, _ := s.repo.GetCredentialByUserID(ctx, user.Id)
		if newCred != nil && newCred.FailedLoginAttempts >= s.cfg.LockoutThreshold {
			lockUntil := time.Now().Add(s.cfg.LockoutDuration)
			s.repo.LockAccount(ctx, user.Id, lockUntil)
			s.logAudit(ctx, &user.Id, string(domain.AuditActionAccountLock), "user", strconv.Itoa(user.Id),
				nil, map[string]interface{}{"locked_until": lockUntil}, req.IPAddress, req.UserAgent)
		}

		s.logFailedLogin(ctx, &user.Id, req.Email, req.IPAddress, req.UserAgent, "invalid password")
		return nil, ErrInvalidCredentials
	}

	// Reset failed attempts on successful password verification
	s.repo.ResetFailedAttempts(ctx, user.Id)

	// Check if MFA is enabled
	mfa, err := s.repo.GetMFAByUserID(ctx, user.Id)
	if err == nil && mfa != nil && mfa.IsEnabled {
		// MFA required - return pending token
		mfaToken := uuid.New().String()
		pendingData := &MFAPendingData{
			UserID:    user.Id,
			Email:     user.Email,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
			ExpiresAt: time.Now().Add(s.cfg.MFATokenTTL),
		}
		if err := s.sessionStore.StoreMFAToken(ctx, mfaToken, pendingData, s.cfg.MFATokenTTL); err != nil {
			return nil, fmt.Errorf("failed to store MFA token: %w", err)
		}

		return &LoginResponse{
			RequiresMFA: true,
			MFAToken:    mfaToken,
		}, nil
	}

	// No MFA - create session directly
	session, err := s.createSession(ctx, user, req.IPAddress, req.UserAgent)
	if err != nil {
		return nil, err
	}

	// Update login stats
	s.repo.UpdateUserLoginStats(ctx, user.Id)

	// Log successful login
	s.logSuccessfulLogin(ctx, user.Id, req.Email, req.IPAddress, req.UserAgent, false)

	return &LoginResponse{
		RequiresMFA: false,
		Session:     session,
		User:        user,
	}, nil
}

// ============================================================
// MFA VERIFICATION
// ============================================================

// VerifyMFARequest contains MFA verification parameters
type VerifyMFARequest struct {
	MFAToken  string
	Code      string
	IPAddress string
	UserAgent string
}

// VerifyMFA verifies a TOTP code and completes login
func (s *AuthService) VerifyMFA(ctx context.Context, req VerifyMFARequest) (*LoginResponse, error) {
	// Get pending MFA data
	pending, err := s.sessionStore.GetMFAToken(ctx, req.MFAToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get MFA token: %w", err)
	}
	if pending == nil {
		return nil, ErrInvalidSession
	}

	// Get MFA config
	mfa, err := s.repo.GetMFAByUserID(ctx, pending.UserID)
	if err != nil || mfa == nil || !mfa.IsEnabled {
		return nil, ErrMFANotEnabled
	}

	// Decrypt secret
	secret, err := s.decryptTOTPSecret(mfa.SecretEncrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt MFA secret: %w", err)
	}

	// Verify TOTP code
	valid := totp.Validate(req.Code, secret)
	if !valid {
		s.logFailedLogin(ctx, &pending.UserID, pending.Email, req.IPAddress, req.UserAgent, "invalid MFA code")
		return nil, ErrInvalidMFACode
	}

	// Delete MFA token
	s.sessionStore.DeleteMFAToken(ctx, req.MFAToken)

	// Get user
	user, err := s.repo.GetUserByEmail(ctx, pending.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Create session
	session, err := s.createSession(ctx, user, req.IPAddress, req.UserAgent)
	if err != nil {
		return nil, err
	}

	// Update login stats
	s.repo.UpdateUserLoginStats(ctx, user.Id)

	// Log successful login with MFA
	s.logSuccessfulLogin(ctx, user.Id, pending.Email, req.IPAddress, req.UserAgent, true)

	return &LoginResponse{
		RequiresMFA: false,
		Session:     session,
		User:        user,
	}, nil
}

// VerifyBackupCode verifies a backup code and completes login
func (s *AuthService) VerifyBackupCode(ctx context.Context, mfaToken, code, ip, userAgent string) (*LoginResponse, error) {
	// Get pending MFA data
	pending, err := s.sessionStore.GetMFAToken(ctx, mfaToken)
	if err != nil || pending == nil {
		return nil, ErrInvalidSession
	}

	// Get backup codes
	codes, err := s.repo.GetUnusedBackupCodes(ctx, pending.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup codes: %w", err)
	}

	// Verify code using each code's unique salt
	var matchedCode *domain.UserMFABackupCode
	for i := range codes {
		// Decode the stored salt
		salt, err := base64.RawStdEncoding.DecodeString(codes[i].Salt)
		if err != nil {
			continue // Skip codes with invalid salt
		}

		// Hash the provided code with the stored salt
		codeHash := s.hashBackupCodeWithSalt(code, salt)
		if subtle.ConstantTimeCompare([]byte(codes[i].CodeHash), []byte(codeHash)) == 1 {
			matchedCode = &codes[i]
			break
		}
	}

	if matchedCode == nil {
		s.logFailedLogin(ctx, &pending.UserID, pending.Email, ip, userAgent, "invalid backup code")
		return nil, ErrInvalidMFACode
	}

	// Mark code as used
	s.repo.UseBackupCode(ctx, matchedCode.ID)

	// Log backup code usage
	s.logAudit(ctx, &pending.UserID, string(domain.AuditActionMFABackupUsed), "backup_code",
		strconv.Itoa(matchedCode.ID), nil, nil, ip, userAgent)

	// Delete MFA token
	s.sessionStore.DeleteMFAToken(ctx, mfaToken)

	// Get user
	user, err := s.repo.GetUserByEmail(ctx, pending.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Create session
	session, err := s.createSession(ctx, user, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Update login stats
	s.repo.UpdateUserLoginStats(ctx, user.Id)

	// Log successful login
	s.logSuccessfulLogin(ctx, user.Id, pending.Email, ip, userAgent, true)

	return &LoginResponse{
		RequiresMFA: false,
		Session:     session,
		User:        user,
	}, nil
}

// ============================================================
// MFA SETUP
// ============================================================

// TOTPSetupResponse contains TOTP setup information
type TOTPSetupResponse struct {
	Secret    string
	QRCodeURL string
}

// SetupTOTP initiates TOTP setup for a user
func (s *AuthService) SetupTOTP(ctx context.Context, userID int, email string) (*TOTPSetupResponse, error) {
	// Check if MFA already enabled
	existing, _ := s.repo.GetMFAByUserID(ctx, userID)
	if existing != nil && existing.IsEnabled {
		return nil, ErrMFAAlreadyEnabled
	}

	// Generate new TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.cfg.TOTPIssuer,
		AccountName: email,
		Period:      30,
		SecretSize:  32,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Encrypt secret
	encryptedSecret, err := s.encryptTOTPSecret(key.Secret())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Store or update MFA record (not enabled yet)
	if existing != nil {
		existing.SecretEncrypted = encryptedSecret
		existing.IsEnabled = false
		existing.VerifiedAt = nil
		if err := s.repo.UpdateMFA(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update MFA: %w", err)
		}
	} else {
		mfa := &domain.UserMFATotp{
			UserID:          userID,
			SecretEncrypted: encryptedSecret,
			IsEnabled:       false,
		}
		if err := s.repo.CreateMFA(ctx, mfa); err != nil {
			return nil, fmt.Errorf("failed to create MFA: %w", err)
		}
	}

	return &TOTPSetupResponse{
		Secret:    key.Secret(),
		QRCodeURL: key.URL(),
	}, nil
}

// ConfirmTOTP confirms TOTP setup with a valid code
func (s *AuthService) ConfirmTOTP(ctx context.Context, userID int, code, ip, userAgent string) error {
	// Get MFA record
	mfa, err := s.repo.GetMFAByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("MFA not set up: %w", err)
	}

	if mfa.IsEnabled {
		return ErrMFAAlreadyEnabled
	}

	// Decrypt secret
	secret, err := s.decryptTOTPSecret(mfa.SecretEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Verify code
	if !totp.Validate(code, secret) {
		return ErrInvalidMFACode
	}

	// Enable MFA
	if err := s.repo.EnableMFA(ctx, userID); err != nil {
		return fmt.Errorf("failed to enable MFA: %w", err)
	}

	// Generate backup codes
	backupCodes, err := s.generateBackupCodes(ctx, userID)
	if err != nil {
		s.logger.Error("failed to generate backup codes", zap.Error(err))
	}

	// Log MFA enable
	s.logAudit(ctx, &userID, string(domain.AuditActionMFAEnable), "user", strconv.Itoa(userID),
		nil, map[string]interface{}{"backup_codes_generated": len(backupCodes)}, ip, userAgent)

	return nil
}

// DisableTOTP disables TOTP for a user
func (s *AuthService) DisableTOTP(ctx context.Context, userID int, code, ip, userAgent string) error {
	// Get MFA record
	mfa, err := s.repo.GetMFAByUserID(ctx, userID)
	if err != nil || mfa == nil {
		return ErrMFANotEnabled
	}

	if !mfa.IsEnabled {
		return ErrMFANotEnabled
	}

	// Decrypt secret
	secret, err := s.decryptTOTPSecret(mfa.SecretEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Verify code
	if !totp.Validate(code, secret) {
		return ErrInvalidMFACode
	}

	// Disable MFA
	if err := s.repo.DisableMFA(ctx, userID); err != nil {
		return fmt.Errorf("failed to disable MFA: %w", err)
	}

	// Delete backup codes
	s.repo.DeleteBackupCodes(ctx, userID)

	// Log MFA disable
	s.logAudit(ctx, &userID, string(domain.AuditActionMFADisable), "user", strconv.Itoa(userID),
		nil, nil, ip, userAgent)

	return nil
}

// ============================================================
// BACKUP CODES
// ============================================================

// GenerateBackupCodes generates new backup codes for a user
func (s *AuthService) GenerateBackupCodes(ctx context.Context, userID int, ip, userAgent string) ([]string, error) {
	// Check MFA is enabled
	mfa, err := s.repo.GetMFAByUserID(ctx, userID)
	if err != nil || mfa == nil || !mfa.IsEnabled {
		return nil, ErrMFANotEnabled
	}

	// Delete existing codes
	s.repo.DeleteBackupCodes(ctx, userID)

	// Generate new codes
	codes, err := s.generateBackupCodes(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Log regeneration
	s.logAudit(ctx, &userID, string(domain.AuditActionMFABackupRegen), "user", strconv.Itoa(userID),
		nil, map[string]interface{}{"codes_generated": len(codes)}, ip, userAgent)

	return codes, nil
}

func (s *AuthService) generateBackupCodes(ctx context.Context, userID int) ([]string, error) {
	codes := make([]string, 10)
	dbCodes := make([]domain.UserMFABackupCode, 10)

	for i := 0; i < 10; i++ {
		code := s.generateRandomCode()
		codes[i] = code

		// Generate unique random salt for each backup code
		salt, saltBase64, err := s.generateBackupCodeSalt()
		if err != nil {
			return nil, fmt.Errorf("failed to generate salt: %w", err)
		}

		dbCodes[i] = domain.UserMFABackupCode{
			UserID:   userID,
			CodeHash: s.hashBackupCodeWithSalt(code, salt),
			Salt:     saltBase64,
		}
	}

	if err := s.repo.CreateBackupCodes(ctx, dbCodes); err != nil {
		return nil, fmt.Errorf("failed to store backup codes: %w", err)
	}

	return codes, nil
}

// ============================================================
// PASSWORD MANAGEMENT
// ============================================================

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID int, currentPass, newPass, ip, userAgent string) error {
	// Get credentials
	cred, err := s.repo.GetCredentialByUserID(ctx, userID)
	if err != nil {
		return ErrCredentialsNotFound
	}

	// Verify current password
	if !s.verifyPassword(currentPass, cred.PasswordHash) {
		return ErrInvalidCredentials
	}

	// Validate new password
	if len(newPass) < s.cfg.PasswordMinLength {
		return ErrPasswordTooWeak
	}

	// Check password history
	if err := s.checkPasswordHistory(ctx, userID, newPass); err != nil {
		return err
	}

	// Hash new password
	newHash, err := s.hashPassword(newPass)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Add current password to history
	s.repo.CreatePasswordHistory(ctx, &domain.PasswordHistory{
		UserID:       userID,
		PasswordHash: cred.PasswordHash,
	})

	// Update credential
	now := time.Now()
	cred.PasswordHash = newHash
	cred.PasswordChangedAt = &now
	cred.MustChangePassword = false

	if err := s.repo.UpdateCredential(ctx, cred); err != nil {
		return fmt.Errorf("failed to update credential: %w", err)
	}

	// Log password change
	s.logAudit(ctx, &userID, string(domain.AuditActionPasswordChange), "user", strconv.Itoa(userID),
		nil, nil, ip, userAgent)

	return nil
}

// SetPassword sets a password for a user (admin/setup)
func (s *AuthService) SetPassword(ctx context.Context, userID int, password string) error {
	// Validate password
	if len(password) < s.cfg.PasswordMinLength {
		return ErrPasswordTooWeak
	}

	// Hash password
	hash, err := s.hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Check if credentials exist
	existing, _ := s.repo.GetCredentialByUserID(ctx, userID)
	now := time.Now()

	if existing != nil {
		existing.PasswordHash = hash
		existing.PasswordChangedAt = &now
		existing.MustChangePassword = true
		return s.repo.UpdateCredential(ctx, existing)
	}

	// Create new credentials
	return s.repo.CreateCredential(ctx, &domain.UserCredential{
		UserID:             userID,
		PasswordHash:       hash,
		PasswordChangedAt:  &now,
		MustChangePassword: true,
	})
}

func (s *AuthService) checkPasswordHistory(ctx context.Context, userID int, newPassword string) error {
	history, err := s.repo.GetPasswordHistory(ctx, userID, s.cfg.PasswordHistoryCount)
	if err != nil {
		return nil // No history, allow
	}

	for _, h := range history {
		if s.verifyPassword(newPassword, h.PasswordHash) {
			return ErrPasswordReused
		}
	}

	return nil
}

// ============================================================
// SESSION MANAGEMENT
// ============================================================

// GetSession retrieves a session by ID
func (s *AuthService) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	return s.sessionStore.Get(ctx, sessionID)
}

// RefreshSession extends a session's expiry
func (s *AuthService) RefreshSession(ctx context.Context, sessionID string) (*SessionData, error) {
	session, err := s.sessionStore.Get(ctx, sessionID)
	if err != nil || session == nil {
		return nil, ErrInvalidSession
	}

	newExpiry := time.Now().Add(s.cfg.SessionTTL)
	if err := s.sessionStore.Refresh(ctx, sessionID, newExpiry); err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	session.ExpiresAt = newExpiry
	session.LastActivityAt = time.Now()

	return session, nil
}

// Logout revokes a session
func (s *AuthService) Logout(ctx context.Context, sessionID, ip, userAgent string) error {
	session, err := s.sessionStore.Get(ctx, sessionID)
	if err != nil || session == nil {
		return nil // Already logged out
	}

	// Delete from Redis
	if err := s.sessionStore.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Revoke in DB
	s.repo.RevokeSession(ctx, sessionID, "user logout")

	// Log
	s.logAudit(ctx, &session.UserID, string(domain.AuditActionSessionRevoke), "session", sessionID,
		nil, map[string]interface{}{"reason": "user logout"}, ip, userAgent)

	return nil
}

// LogoutAll revokes all sessions for a user
func (s *AuthService) LogoutAll(ctx context.Context, userID int, ip, userAgent string) error {
	// Delete all sessions from Redis
	if err := s.sessionStore.DeleteAllUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	// Revoke all in DB
	s.repo.RevokeAllUserSessions(ctx, userID, "logout all")

	// Log
	s.logAudit(ctx, &userID, string(domain.AuditActionLogoutAll), "user", strconv.Itoa(userID),
		nil, nil, ip, userAgent)

	return nil
}

// GetActiveSessions returns all active sessions for a user
func (s *AuthService) GetActiveSessions(ctx context.Context, userID int) ([]SessionData, error) {
	sessionIDs, err := s.sessionStore.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sessions []SessionData
	for _, id := range sessionIDs {
		session, err := s.sessionStore.Get(ctx, id)
		if err == nil && session != nil {
			sessions = append(sessions, *session)
		}
	}

	return sessions, nil
}

func (s *AuthService) createSession(ctx context.Context, user *domain.User, ip, userAgent string) (*SessionData, error) {
	sessionID := uuid.New().String()
	now := time.Now()

	session := &SessionData{
		SessionID:      sessionID,
		UserID:         user.Id,
		Email:          user.Email,
		IPAddress:      ip,
		UserAgent:      userAgent,
		CreatedAt:      now,
		ExpiresAt:      now.Add(s.cfg.SessionTTL),
		LastActivityAt: now,
	}

	// Store in Redis
	if err := s.sessionStore.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Store in DB for audit
	dbSession := &domain.Session{
		ID:             sessionID,
		UserID:         user.Id,
		IPAddress:      ip,
		UserAgent:      userAgent,
		ExpiresAt:      session.ExpiresAt,
		LastActivityAt: now,
	}
	s.repo.CreateSession(ctx, dbSession)

	return session, nil
}

// ============================================================
// USER STATUS MANAGEMENT
// ============================================================

// UpdateUserStatus updates a user's status
func (s *AuthService) UpdateUserStatus(ctx context.Context, userID int, status, reason string, changedBy int, ip, userAgent string) error {
	// Validate status
	userStatus := domain.UserStatus(status)
	if !userStatus.IsValid() {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Get current status for audit
	user, err := s.repo.GetUserByEmail(ctx, "") // This needs fixing - we need GetUserByID
	if err != nil {
		s.logger.Warn("could not get user for audit", zap.Int("user_id", userID))
	}

	oldStatus := ""
	if user != nil {
		oldStatus = user.Status
	}

	// Update status
	if err := s.repo.UpdateUserStatus(ctx, userID, status, reason, changedBy); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// If locked/suspended, revoke all sessions
	if status == string(domain.UserStatusLocked) || status == string(domain.UserStatusSuspended) {
		s.sessionStore.DeleteAllUserSessions(ctx, userID)
		s.repo.RevokeAllUserSessions(ctx, userID, "status change: "+status)
	}

	// Log audit
	s.logAudit(ctx, &changedBy, string(domain.AuditActionStatusChange), "user", strconv.Itoa(userID),
		map[string]interface{}{"status": oldStatus},
		map[string]interface{}{"status": status, "reason": reason},
		ip, userAgent)

	return nil
}

// UnlockAccount unlocks a locked account
func (s *AuthService) UnlockAccount(ctx context.Context, userID int, unlockedBy int, ip, userAgent string) error {
	if err := s.repo.UnlockAccount(ctx, userID); err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	// Log
	s.logAudit(ctx, &unlockedBy, string(domain.AuditActionAccountUnlock), "user", strconv.Itoa(userID),
		nil, nil, ip, userAgent)

	return nil
}

// ============================================================
// AUDIT & HISTORY
// ============================================================

// GetLoginHistory returns login history for a user
func (s *AuthService) GetLoginHistory(ctx context.Context, userID int, limit int) ([]domain.LoginHistory, error) {
	return s.repo.GetLoginHistory(ctx, userID, limit)
}

// GetSecurityAudit returns security audit trail for a user
func (s *AuthService) GetSecurityAudit(ctx context.Context, userID int, limit int) ([]domain.SecurityAuditTrail, error) {
	return s.repo.GetAuditTrail(ctx, userID, limit)
}

// GetMFAStatus returns MFA status for a user
func (s *AuthService) GetMFAStatus(ctx context.Context, userID int) (enabled bool, hasBackupCodes bool, err error) {
	mfa, err := s.repo.GetMFAByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, nil
		}
		return false, false, err
	}

	if !mfa.IsEnabled {
		return false, false, nil
	}

	codes, _ := s.repo.GetUnusedBackupCodes(ctx, userID)
	return true, len(codes) > 0, nil
}

// ============================================================
// HELPER METHODS
// ============================================================

func (s *AuthService) hashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Encode as: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads, b64Salt, b64Hash)

	return encoded, nil
}

func (s *AuthService) verifyPassword(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	var version int
	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false
	}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	return subtle.ConstantTimeCompare(hash, comparisonHash) == 1
}

func (s *AuthService) encryptTOTPSecret(secret string) (string, error) {
	key := []byte(s.cfg.EncryptionKey)
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(secret), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *AuthService) decryptTOTPSecret(encrypted string) (string, error) {
	key := []byte(s.cfg.EncryptionKey)
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (s *AuthService) generateRandomCode() string {
	b := make([]byte, 5)
	rand.Read(b)
	return strings.ToUpper(base32.StdEncoding.EncodeToString(b)[:8])
}

// hashBackupCodeWithSalt hashes a backup code using the provided salt
func (s *AuthService) hashBackupCodeWithSalt(code string, salt []byte) string {
	hash := argon2.IDKey([]byte(code), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash)
}

// generateBackupCodeSalt generates a random salt for backup code hashing
func (s *AuthService) generateBackupCodeSalt() ([]byte, string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, "", err
	}
	return salt, base64.RawStdEncoding.EncodeToString(salt), nil
}

func (s *AuthService) logFailedLogin(ctx context.Context, userID *int, email, ip, userAgent, reason string) {
	history := &domain.LoginHistory{
		UserID:        userID,
		Email:         email,
		IPAddress:     ip,
		UserAgent:     userAgent,
		LoginMethod:   "local",
		Success:       false,
		FailureReason: reason,
	}
	s.repo.CreateLoginHistory(ctx, history)
}

func (s *AuthService) logSuccessfulLogin(ctx context.Context, userID int, email, ip, userAgent string, mfaUsed bool) {
	history := &domain.LoginHistory{
		UserID:      &userID,
		Email:       email,
		IPAddress:   ip,
		UserAgent:   userAgent,
		LoginMethod: "local",
		Success:     true,
		MFAUsed:     mfaUsed,
	}
	s.repo.CreateLoginHistory(ctx, history)

	// Also log audit
	s.logAudit(ctx, &userID, string(domain.AuditActionLoginSuccess), "user", strconv.Itoa(userID),
		nil, map[string]interface{}{"mfa_used": mfaUsed}, ip, userAgent)
}

func (s *AuthService) logAudit(ctx context.Context, userID *int, action, targetType, targetID string, oldValue, newValue interface{}, ip, userAgent string) {
	var oldJSON, newJSON string
	if oldValue != nil {
		if b, err := json.Marshal(oldValue); err == nil {
			oldJSON = string(b)
		}
	}
	if newValue != nil {
		if b, err := json.Marshal(newValue); err == nil {
			newJSON = string(b)
		}
	}

	audit := &domain.SecurityAuditTrail{
		UserID:     userID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		OldValue:   oldJSON,
		NewValue:   newJSON,
		IPAddress:  ip,
		UserAgent:  userAgent,
	}
	s.repo.CreateAuditTrail(ctx, audit)
}
