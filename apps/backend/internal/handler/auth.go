package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/middleware"
	"github.com/reche13/habitum/internal/model/user"
	"github.com/reche13/habitum/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req user.LoginRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	authResp, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, authResp)
}

// Signup handles POST /api/v1/auth/signup
func (h *AuthHandler) Signup(c echo.Context) error {
	var req user.SignupRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	authResp, err := h.authService.Signup(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, authResp)
}

// GoogleAuth handles POST /api/v1/auth/google
func (h *AuthHandler) GoogleAuth(c echo.Context) error {
	var req user.GoogleAuthRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	authResp, err := h.authService.GoogleAuth(c.Request().Context(), req.Token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, authResp)
}

// VerifyEmail handles POST /api/v1/auth/verify-email
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var req user.VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	if err := h.authService.VerifyEmail(c.Request().Context(), req.Token); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email verified successfully"})
}

// ResendVerification handles POST /api/v1/auth/resend-verification
func (h *AuthHandler) ResendVerification(c echo.Context) error {
	var req user.ResendVerificationRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	if err := h.authService.ResendVerificationEmail(c.Request().Context(), req.Email); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Verification email sent"})
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req user.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	if err := h.authService.ForgotPassword(c.Request().Context(), req.Email); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset email sent"})
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req user.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	if err := h.authService.ResetPassword(c.Request().Context(), req.Token, req.NewPassword); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req user.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&req); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	authResp, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, authResp)
}

// TestAccountLogin handles POST /api/v1/auth/test-account
func (h *AuthHandler) TestAccountLogin(c echo.Context) error {
	authResp, err := h.authService.TestAccountLogin(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, authResp)
}

// Logout handles POST /api/v1/auth/logout
// This is mainly for client-side token removal, but we can add token blacklisting here if needed
func (h *AuthHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

