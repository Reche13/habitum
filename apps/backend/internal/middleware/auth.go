package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/service"
)

type AuthContext struct {
	UserID uuid.UUID
	Email  string
}

// AuthMiddleware validates JWT token and adds user context
func AuthMiddleware(jwtService *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errs.NewUnauthorizedError("missing authorization header")
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return errs.NewUnauthorizedError("invalid authorization header format")
			}

			token := parts[1]

			// Validate token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return errs.NewUnauthorizedError("invalid or expired token")
			}

			// Check token type (should be access token)
			if claims.Type != "access" {
				return errs.NewUnauthorizedError("invalid token type")
			}

			// Parse user ID
			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				return errs.NewUnauthorizedError("invalid user ID in token")
			}

			// Add user context to request
			c.Set("user_id", userID)
			c.Set("user_email", claims.Email)
			c.Set("auth_context", &AuthContext{
				UserID: userID,
				Email:  claims.Email,
			})

			return next(c)
		}
	}
}

// GetUserID extracts user ID from context
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errs.NewUnauthorizedError("user ID not found in context")
	}
	return userID, nil
}

// GetAuthContext extracts auth context from request
func GetAuthContext(c echo.Context) (*AuthContext, error) {
	ctx, ok := c.Get("auth_context").(*AuthContext)
	if !ok {
		return nil, errs.NewUnauthorizedError("auth context not found")
	}
	return ctx, nil
}

