// Package service provides implementation for service
//
// File: registration_service.go
// Description: Registration service for user signup, email verification, and password reset
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"templatev25/internal/config"
	"templatev25/internal/domain"
	"templatev25/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Registration error definitions
var (
	ErrEmailAlreadyExists       = errors.New("email already registered")
	ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
	ErrInvalidResetToken        = errors.New("invalid or expired password reset token")
	ErrUserAlreadyVerified      = errors.New("user is already verified")
	ErrPasswordMismatch         = errors.New("passwords do not match")
)

// RegistrationService handles user registration, email verification, and password reset
type RegistrationService struct {
	authRepo    repository.AuthRepository
	userRepo    repository.UserRepository
	regRepo     repository.RegistrationRepository
	authService *AuthService
	cfg         *config.LocalAuthConfig
	logger      *zap.Logger
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	regRepo repository.RegistrationRepository,
	authService *AuthService,
	cfg *config.LocalAuthConfig,
	logger *zap.Logger,
) *RegistrationService {
	return &RegistrationService{
		authRepo:    authRepo,
		userRepo:    userRepo,
		regRepo:     regRepo,
		authService: authService,
		cfg:         cfg,
		logger:      logger,
	}
}

// ============================================================
// REGISTRATION
// ============================================================

// RegisterRequest contains registration parameters
type RegistrationRequest struct {
	Email           string
	Password        string
	ConfirmPassword string
	FirstName       string
	LastName        string
	IPAddress       string
	UserAgent       string
}

// RegisterResponse contains registration result
type RegistrationResponse struct {
	UserID           int
	Email            string
	VerificationSent bool
	Message          string
}

// Register creates a new user account
func (s *RegistrationService) Register(ctx context.Context, req RegistrationRequest) (*RegistrationResponse, error) {
	// Validate password match
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	// Validate password strength
	if len(req.Password) < s.cfg.PasswordMinLength {
		return nil, ErrPasswordTooWeak
	}

	// Check if email already exists
	exists, err := s.regRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Create user with pending_verification status
	user := &domain.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    string(domain.UserStatusPendingVerification),
	}

	if err := s.regRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Set password
	if err := s.authService.SetPassword(ctx, user.Id, req.Password); err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	// Generate verification token
	token, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	verificationToken := &domain.EmailVerificationToken{
		UserID:    user.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
	}

	if err := s.regRepo.CreateEmailVerificationToken(ctx, verificationToken); err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// TODO: Send verification email
	// s.emailService.SendVerificationEmail(user.Email, token)

	s.logger.Info("user registered",
		zap.Int("user_id", user.Id),
		zap.String("email", user.Email),
	)

	return &RegistrationResponse{
		UserID:           user.Id,
		Email:            user.Email,
		VerificationSent: true,
		Message:          "Registration successful. Please check your email to verify your account.",
	}, nil
}

// ============================================================
// EMAIL VERIFICATION
// ============================================================

// VerifyEmail verifies a user's email address
func (s *RegistrationService) VerifyEmail(ctx context.Context, tokenStr string) error {
	// Get token
	token, err := s.regRepo.GetEmailVerificationToken(ctx, tokenStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidVerificationToken
		}
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Check if token is valid
	if token.IsExpired() || token.IsUsed() {
		return ErrInvalidVerificationToken
	}

	// Mark token as used
	if err := s.regRepo.MarkEmailVerificationTokenUsed(ctx, token.ID); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	// Update user as verified
	if err := s.regRepo.UpdateUserEmailVerified(ctx, token.UserID); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Update user status to active
	if err := s.authRepo.UpdateUserStatus(ctx, token.UserID, string(domain.UserStatusActive), "email verified", 0); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	s.logger.Info("email verified",
		zap.Int("user_id", token.UserID),
	)

	return nil
}

// ResendVerificationEmail resends the verification email
func (s *RegistrationService) ResendVerificationEmail(ctx context.Context, email string) error {
	// Get user
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if already verified
	if user.Status == string(domain.UserStatusActive) {
		return ErrUserAlreadyVerified
	}

	// Delete existing tokens
	s.regRepo.DeleteUserEmailVerificationTokens(ctx, user.Id)

	// Generate new token
	token, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	verificationToken := &domain.EmailVerificationToken{
		UserID:    user.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.regRepo.CreateEmailVerificationToken(ctx, verificationToken); err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	// TODO: Send verification email
	// s.emailService.SendVerificationEmail(user.Email, token)

	return nil
}

// ============================================================
// PASSWORD RESET
// ============================================================

// ForgotPassword initiates the password reset process
func (s *RegistrationService) ForgotPassword(ctx context.Context, email string) error {
	// Get user
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists - always return success
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete existing tokens
	s.regRepo.DeleteUserPasswordResetTokens(ctx, user.Id)

	// Generate new token
	token, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	resetToken := &domain.PasswordResetToken{
		UserID:    user.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	if err := s.regRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send password reset email
	// s.emailService.SendPasswordResetEmail(user.Email, token)

	s.logger.Info("password reset requested",
		zap.Int("user_id", user.Id),
		zap.String("email", user.Email),
	)

	return nil
}

// ResetPassword resets the user's password using a valid token
func (s *RegistrationService) ResetPassword(ctx context.Context, tokenStr, newPassword, confirmPassword string) error {
	// Validate password match
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}

	// Validate password strength
	if len(newPassword) < s.cfg.PasswordMinLength {
		return ErrPasswordTooWeak
	}

	// Get token
	token, err := s.regRepo.GetPasswordResetToken(ctx, tokenStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Check if token is valid
	if token.IsExpired() || token.IsUsed() {
		return ErrInvalidResetToken
	}

	// Mark token as used
	if err := s.regRepo.MarkPasswordResetTokenUsed(ctx, token.ID); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	// Set new password
	if err := s.authService.SetPassword(ctx, token.UserID, newPassword); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	// Delete all password reset tokens for this user
	s.regRepo.DeleteUserPasswordResetTokens(ctx, token.UserID)

	// Revoke all sessions for security
	s.authService.LogoutAll(ctx, token.UserID, "", "password reset")

	s.logger.Info("password reset successful",
		zap.Int("user_id", token.UserID),
	)

	return nil
}

// ============================================================
// HELPER METHODS
// ============================================================

// generateSecureToken generates a cryptographically secure token
func (s *RegistrationService) generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
