package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/handler"
)

func registerAuthRoutes(auth *echo.Group, h *handler.Handlers) {
	auth.POST("/login", h.Auth.Login)
	auth.POST("/signup", h.Auth.Signup)
	auth.POST("/google", h.Auth.GoogleAuth)
	auth.POST("/verify-email", h.Auth.VerifyEmail)
	auth.POST("/resend-verification", h.Auth.ResendVerification)
	auth.POST("/forgot-password", h.Auth.ForgotPassword)
	auth.POST("/reset-password", h.Auth.ResetPassword)
	auth.POST("/refresh", h.Auth.RefreshToken)
	auth.POST("/test-account", h.Auth.TestAccountLogin)
	auth.POST("/logout", h.Auth.Logout)
}

