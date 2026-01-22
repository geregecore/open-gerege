// Package handlers provides implementation for handlers
//
// File: registration_handler.go
// Description: Handler for user registration, email verification, and password reset
package handlers

import (
	"errors"

	"templatev25/internal/http/dto"
	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/resp"
	"github.com/gofiber/fiber/v2"
)

// RegistrationHandler handles registration-related endpoints
type RegistrationHandler struct {
	registrationService *service.RegistrationService
}

// NewRegistrationHandler creates a new registration handler
func NewRegistrationHandler(registrationService *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		registrationService: registrationService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email and password
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterRequest true "Registration data"
// @Success      201 {object} dto.RegisterResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse "Email already exists"
// @Router       /auth/local/register [post]
func (h *RegistrationHandler) Register(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.RegisterRequest](c)
	if !ok {
		return nil
	}

	regReq := service.RegistrationRequest{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		IPAddress:       c.IP(),
		UserAgent:       c.Get("User-Agent"),
	}

	result, err := h.registrationService.Register(c.UserContext(), regReq)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": "email already registered",
			})
		case errors.Is(err, service.ErrPasswordMismatch):
			return resp.BadRequest(c, "passwords do not match", nil)
		case errors.Is(err, service.ErrPasswordTooWeak):
			return resp.BadRequest(c, "password does not meet requirements", nil)
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(dto.RegisterResponse{
		UserID:           result.UserID,
		Email:            result.Email,
		VerificationSent: result.VerificationSent,
		Message:          result.Message,
	})
}

// VerifyEmail godoc
// @Summary      Verify email address
// @Description  Verify a user's email address using the verification token
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        body body dto.VerifyEmailRequest true "Verification token"
// @Success      200 {object} dto.VerifyEmailResponse
// @Failure      400 {object} dto.ErrorResponse "Invalid or expired token"
// @Router       /auth/local/verify-email [post]
func (h *RegistrationHandler) VerifyEmail(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.VerifyEmailRequest](c)
	if !ok {
		return nil
	}

	err := h.registrationService.VerifyEmail(c.UserContext(), req.Token)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidVerificationToken):
			return resp.BadRequest(c, "invalid or expired verification token", nil)
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	return resp.OK(c, dto.VerifyEmailResponse{
		Success: true,
		Message: "Email verified successfully. You can now log in.",
	})
}

// ResendVerification godoc
// @Summary      Resend verification email
// @Description  Resend the email verification link to the user's email
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        body body dto.ResendVerificationRequest true "Email address"
// @Success      200 {object} dto.GenericResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /auth/local/resend-verification [post]
func (h *RegistrationHandler) ResendVerification(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ResendVerificationRequest](c)
	if !ok {
		return nil
	}

	err := h.registrationService.ResendVerificationEmail(c.UserContext(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyVerified):
			return resp.BadRequest(c, "user is already verified", nil)
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	// Always return success to prevent email enumeration
	return resp.OK(c, dto.GenericResponse{
		Success: true,
		Message: "If an account exists with this email, a verification link has been sent.",
	})
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Send a password reset link to the user's email
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        body body dto.ForgotPasswordRequest true "Email address"
// @Success      200 {object} dto.GenericResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /auth/local/forgot-password [post]
func (h *RegistrationHandler) ForgotPassword(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ForgotPasswordRequest](c)
	if !ok {
		return nil
	}

	err := h.registrationService.ForgotPassword(c.UserContext(), req.Email)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	// Always return success to prevent email enumeration
	return resp.OK(c, dto.GenericResponse{
		Success: true,
		Message: "If an account exists with this email, a password reset link has been sent.",
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset the user's password using a valid reset token
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        body body dto.ResetPasswordConfirmRequest true "Reset token and new password"
// @Success      200 {object} dto.GenericResponse
// @Failure      400 {object} dto.ErrorResponse "Invalid token or password"
// @Router       /auth/local/reset-password [post]
func (h *RegistrationHandler) ResetPassword(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ResetPasswordConfirmRequest](c)
	if !ok {
		return nil
	}

	err := h.registrationService.ResetPassword(c.UserContext(), req.Token, req.Password, req.ConfirmPassword)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidResetToken):
			return resp.BadRequest(c, "invalid or expired reset token", nil)
		case errors.Is(err, service.ErrPasswordMismatch):
			return resp.BadRequest(c, "passwords do not match", nil)
		case errors.Is(err, service.ErrPasswordTooWeak):
			return resp.BadRequest(c, "password does not meet requirements", nil)
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	return resp.OK(c, dto.GenericResponse{
		Success: true,
		Message: "Password reset successfully. You can now log in with your new password.",
	})
}
