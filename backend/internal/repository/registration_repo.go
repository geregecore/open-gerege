// Package repository provides implementation for repository
//
// File: registration_repo.go
// Description: Repository for registration-related entities (email verification, password reset)
package repository

import (
	"context"
	"time"

	"templatev25/internal/domain"

	"gorm.io/gorm"
)

// RegistrationRepository defines additional repository methods for registration
type RegistrationRepository interface {
	// Email verification
	CreateEmailVerificationToken(ctx context.Context, token *domain.EmailVerificationToken) error
	GetEmailVerificationToken(ctx context.Context, token string) (*domain.EmailVerificationToken, error)
	MarkEmailVerificationTokenUsed(ctx context.Context, tokenID int) error
	DeleteUserEmailVerificationTokens(ctx context.Context, userID int) error

	// Password reset
	CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenID int) error
	DeleteUserPasswordResetTokens(ctx context.Context, userID int) error

	// User management
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUserEmailVerified(ctx context.Context, userID int) error
	GetUserByID(ctx context.Context, userID int) (*domain.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

type registrationRepository struct {
	db *gorm.DB
}

// NewRegistrationRepository creates a new registration repository instance
func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

// ============================================================
// EMAIL VERIFICATION TOKENS
// ============================================================

func (r *registrationRepository) CreateEmailVerificationToken(ctx context.Context, token *domain.EmailVerificationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *registrationRepository) GetEmailVerificationToken(ctx context.Context, tokenStr string) (*domain.EmailVerificationToken, error) {
	var token domain.EmailVerificationToken
	err := r.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *registrationRepository) MarkEmailVerificationTokenUsed(ctx context.Context, tokenID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.EmailVerificationToken{}).
		Where("id = ?", tokenID).
		Update("used_at", now).Error
}

func (r *registrationRepository) DeleteUserEmailVerificationTokens(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.EmailVerificationToken{}).Error
}

// ============================================================
// PASSWORD RESET TOKENS
// ============================================================

func (r *registrationRepository) CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *registrationRepository) GetPasswordResetToken(ctx context.Context, tokenStr string) (*domain.PasswordResetToken, error) {
	var token domain.PasswordResetToken
	err := r.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *registrationRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Update("used_at", now).Error
}

func (r *registrationRepository) DeleteUserPasswordResetTokens(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.PasswordResetToken{}).Error
}

// ============================================================
// USER MANAGEMENT
// ============================================================

func (r *registrationRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *registrationRepository) UpdateUserEmailVerified(ctx context.Context, userID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email_verified":    true,
			"email_verified_at": now,
		}).Error
}

func (r *registrationRepository) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_date IS NULL", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *registrationRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("email = ? AND deleted_date IS NULL", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
