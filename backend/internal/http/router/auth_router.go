// Package router provides implementation for router
//
// File: auth_router.go
// Description: Authentication routes implementation
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package router

import (
	"time"

	"templatev25/internal/app"
	"templatev25/internal/http/handlers"
	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// MapAuthRoutes нь authentication-тай холбоотой route-уудыг бүртгэнэ.
//
// Routes:
//   SSO Routes:
//   - GET  /auth/login      → SSO redirect
//   - GET  /auth/callback   → OAuth2 callback
//   - POST /auth/logout     → Logout
//   - POST /auth/google/login → Google OAuth
//   - GET  /auth/verify     → Token verification
//   - POST /auth/org/change → Change organization (protected)
//
//   Local Auth Routes:
//   - POST /auth/local/login        → Local login with email/password
//   - POST /auth/local/verify-mfa   → Verify MFA code
//   - POST /auth/local/verify-backup → Verify backup code
//   - POST /auth/local/logout       → Local logout (protected)
//   - POST /auth/local/logout-all   → Logout all sessions (protected)
//   - POST /auth/local/refresh      → Refresh session (protected)
//
// Security:
//   - AuthRateLimiter: 5 req/min per IP for login/callback (brute force protection)
//   - StrictRateLimiter: 3 req/5min for sensitive operations
func MapAuthRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// ------------------------------------------------------------
	// SSO AUTH ROUTES
	// ------------------------------------------------------------
	// Authentication-тай холбоотой endpoint-ууд.
	// Timeout: 5 секунд (SSO response хүлээх)
	v1.Group("/auth", middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		handler := handlers.NewAuthHandler(d)

		// Rate limiter for auth endpoints (brute force protection)
		authLimiter := middleware.AuthRateLimiter()

		// SSO login redirect
		// GET /auth/login → Redirect to SSO login page
		// Rate limited: 5 req/min per IP
		router.Get("/login", authLimiter, handler.InitDirection)

		// OAuth2 callback (SSO-оос буцаж ирэх)
		// GET /auth/callback?code=xxx → Exchange code for token
		// Rate limited: 5 req/min per IP
		router.Get("/callback", authLimiter, handler.OAuthCallback)

		// Logout
		// POST /auth/logout → Clear session, redirect to SSO logout
		router.Post("/logout", handler.Logout)

		// Google OAuth login
		// POST /auth/google/login → Google OAuth flow
		// Rate limited: 5 req/min per IP
		router.Post("/google/login", authLimiter, handler.GoogleLogin)

		// Token verification
		// GET /auth/verify → Check if current session is valid
		router.Get("/verify", handler.AuthVerify)

		// Change organization (protected)
		// POST /auth/org/change → Switch to different organization
		// Strict rate limit: 3 req/5min (sensitive operation)
		router.Post("/org/change", requireAuth, middleware.StrictRateLimiter(), handler.ChangeOrganization)
	})

	// ------------------------------------------------------------
	// LOCAL AUTH ROUTES
	// ------------------------------------------------------------
	// Local authentication with email/password + MFA support
	// Session auth middleware for protected routes
	// Use adapter to bridge service.SessionStore to middleware.SessionStore interface
	sessionStoreAdapter := NewSessionStoreAdapter(d.Service.SessionStore)
	sessionAuth := middleware.SessionAuth(sessionStoreAdapter)

	v1.Group("/auth/local", middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		localAuthHandler := handlers.NewLocalAuthHandler(d.Service.Auth)

		// Rate limiter for auth endpoints (brute force protection)
		authLimiter := middleware.AuthRateLimiter()
		strictLimiter := middleware.StrictRateLimiter()

		// Local login with email/password
		// POST /auth/local/login → Authenticate with email/password
		// Returns session token or MFA token if MFA is enabled
		router.Post("/login", authLimiter, localAuthHandler.Login)

		// MFA verification
		// POST /auth/local/verify-mfa → Verify TOTP code
		router.Post("/verify-mfa", authLimiter, localAuthHandler.VerifyMFA)

		// Backup code verification
		// POST /auth/local/verify-backup → Verify backup code
		router.Post("/verify-backup", authLimiter, localAuthHandler.VerifyBackupCode)

		// Logout (protected by session auth)
		// POST /auth/local/logout → Revoke current session
		router.Post("/logout", sessionAuth, localAuthHandler.Logout)

		// Logout all sessions (protected by session auth)
		// POST /auth/local/logout-all → Revoke all user sessions
		router.Post("/logout-all", sessionAuth, localAuthHandler.LogoutAll)

		// Refresh session (protected by session auth)
		// POST /auth/local/refresh → Extend session expiry
		router.Post("/refresh", sessionAuth, localAuthHandler.RefreshSession)

		// ------------------------------------------------------------
		// REGISTRATION ROUTES (Public)
		// ------------------------------------------------------------
		// Registration handler - only create if registration service exists
		if d.Service.Registration != nil {
			registrationHandler := handlers.NewRegistrationHandler(d.Service.Registration)

			// User registration
			// POST /auth/local/register → Create new user account
			// Rate limited: Strict (3 req/5min) to prevent abuse
			router.Post("/register", strictLimiter, registrationHandler.Register)

			// Email verification
			// POST /auth/local/verify-email → Verify email with token
			router.Post("/verify-email", authLimiter, registrationHandler.VerifyEmail)

			// Resend verification email
			// POST /auth/local/resend-verification → Resend verification email
			// Rate limited: Strict (3 req/5min) to prevent spam
			router.Post("/resend-verification", strictLimiter, registrationHandler.ResendVerification)

			// Password reset flow
			// POST /auth/local/forgot-password → Request password reset email
			// Rate limited: Strict (3 req/5min) to prevent abuse
			router.Post("/forgot-password", strictLimiter, registrationHandler.ForgotPassword)

			// POST /auth/local/reset-password → Reset password with token
			router.Post("/reset-password", authLimiter, registrationHandler.ResetPassword)
		}
	})

	// ------------------------------------------------------------
	// VERIFY ROUTES
	// ------------------------------------------------------------
	// Баталгаажуулалт (DAN, email, phone).
	// Strict rate limiting: 3 req/5min (OTP/verification abuse prevention)
	v1.Group("/verify", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewVerifyHandler(d)
		strictLimiter := middleware.StrictRateLimiter()

		// DAN verification
		router.Get("/dan", h.Dan)

		// Email verification (rate limited - OTP abuse prevention)
		router.Post("/email", strictLimiter, h.Email)
		router.Post("/email/confirm", strictLimiter, h.EmailConfirm)

		// Phone verification (rate limited - SMS abuse prevention)
		router.Post("/phone", strictLimiter, h.Phone)
		router.Post("/phone/confirm", strictLimiter, h.PhoneConfirm)
	})
}
