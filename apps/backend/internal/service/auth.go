package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/model/user"
	"github.com/reche13/habitum/internal/repository"
	"github.com/rs/zerolog"
)

type AuthService struct {
	*BaseService
	userRepo     *repository.UserRepository
	jwtService   *JWTService
	emailService *EmailService
	oauthService *OAuthService
	logger       zerolog.Logger
	testEmail    string
	testPassword string
}

func NewAuthService(
	userRepo *repository.UserRepository,
	jwtService *JWTService,
	emailService *EmailService,
	oauthService *OAuthService,
	logger zerolog.Logger,
	testEmail, testPassword string,
) *AuthService {
	return &AuthService{
		BaseService: &BaseService{
			resourceName: "auth",
		},
		userRepo:     userRepo,
		jwtService:   jwtService,
		emailService: emailService,
		oauthService: oauthService,
		logger:       logger,
		testEmail:    testEmail,
		testPassword: testPassword,
	}
}

// Login handles email/password login
func (s *AuthService) Login(ctx context.Context, email, password string) (*user.AuthResponse, error) {
	// Find user by email
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errs.NewUnauthorizedError("invalid email or password")
	}

	// Check if user has a password (not OAuth-only user)
	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return nil, errs.NewUnauthorizedError("invalid email or password")
	}

	// Verify password
	if !VerifyPassword(password, *u.PasswordHash) {
		return nil, errs.NewUnauthorizedError("invalid email or password")
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, u.ID)

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user.AuthResponse{
		User: &user.UserResponse{
			ID:            u.ID.String(),
			Name:          u.Name,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			OAuthProvider: u.OAuthProvider,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}, nil
}

// Signup handles email/password signup
func (s *AuthService) Signup(ctx context.Context, name, email, password string) (*user.AuthResponse, error) {
	// Validate password strength
	if err := ValidatePasswordStrength(password); err != nil {
		return nil, errs.NewBadRequestError(err.Error())
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errs.NewConflictError("user with this email already exists")
	}

	// Hash password
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification token
	verificationToken, err := GenerateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	expiresAt := GetEmailVerificationExpiry()

	// Create user
	u, err := s.userRepo.CreateWithPassword(ctx, name, email, passwordHash)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Set verification token
	err = s.userRepo.UpdateEmailVerification(ctx, u.ID, false, &verificationToken, &expiresAt)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Send verification email
	if s.emailService != nil {
		if err := s.emailService.SendVerificationEmail(email, name, verificationToken); err != nil {
			s.logger.Warn().Err(err).Str("email", email).Msg("failed to send verification email")
			// Don't fail signup if email fails
		}
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user.AuthResponse{
		User: &user.UserResponse{
			ID:            u.ID.String(),
			Name:          u.Name,
			Email:         u.Email,
			EmailVerified: false,
			OAuthProvider: nil,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}

// GoogleAuth handles Google OAuth login/signup
func (s *AuthService) GoogleAuth(ctx context.Context, idToken string) (*user.AuthResponse, error) {
	// Check if OAuth service is configured
	if s.oauthService == nil {
		return nil, errs.NewBadRequestError("Google OAuth is not configured")
	}

	// Verify Google token
	googleUser, err := s.oauthService.VerifyGoogleToken(ctx, idToken)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to verify Google token")
		return nil, errs.NewUnauthorizedError("invalid Google token: " + err.Error())
	}

	// Check if user exists by OAuth provider
	u, err := s.userRepo.GetByOAuthProvider(ctx, "google", googleUser.ID)
	if err != nil {
		// User doesn't exist, check if email exists
		existingUser, _ := s.userRepo.GetByEmail(ctx, googleUser.Email)
		if existingUser != nil {
			return nil, errs.NewConflictError("user with this email already exists")
		}

		// Create new OAuth user
		u, err = s.userRepo.CreateOAuthUser(ctx, googleUser.Name, googleUser.Email, "google", googleUser.ID)
		if err != nil {
			return nil, s.wrapError(err)
		}
	} else {
		// Update last login
		_ = s.userRepo.UpdateLastLogin(ctx, u.ID)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user.AuthResponse{
		User: &user.UserResponse{
			ID:            u.ID.String(),
			Name:          u.Name,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			OAuthProvider: u.OAuthProvider,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}

// VerifyEmail verifies user email with token
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	u, err := s.userRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		return errs.NewBadRequestError("invalid or expired verification token")
	}

	// Check if token is expired
	if u.EmailVerificationExpiresAt != nil && u.EmailVerificationExpiresAt.Before(time.Now()) {
		return errs.NewBadRequestError("verification token has expired")
	}

	// Verify email
	err = s.userRepo.UpdateEmailVerification(ctx, u.ID, true, nil, nil)
	if err != nil {
		return s.wrapError(err)
	}

	return nil
}

// ResendVerificationEmail resends verification email
func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) error {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not
		return nil
	}

	if u.EmailVerified {
		return nil
	}

	// Generate new token
	verificationToken, err := GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	expiresAt := GetEmailVerificationExpiry()

	// Update token
	err = s.userRepo.UpdateEmailVerification(ctx, u.ID, false, &verificationToken, &expiresAt)
	if err != nil {
		return s.wrapError(err)
	}

	// Send email
	if s.emailService != nil {
		return s.emailService.SendVerificationEmail(u.Email, u.Name, verificationToken)
	}

	return nil
}

// ForgotPassword sends password reset email
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not
		return nil
	}

	// Generate reset token
	resetToken, err := GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expiresAt := GetPasswordResetExpiry()

	// Update token
	err = s.userRepo.UpdatePasswordResetToken(ctx, u.ID, &resetToken, &expiresAt)
	if err != nil {
		return s.wrapError(err)
	}

	// Send email
	if s.emailService != nil {
		return s.emailService.SendPasswordResetEmail(u.Email, u.Name, resetToken)
	}

	return nil
}

// ResetPassword resets password with token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate password strength
	if err := ValidatePasswordStrength(newPassword); err != nil {
		return errs.NewBadRequestError(err.Error())
	}

	u, err := s.userRepo.GetByPasswordResetToken(ctx, token)
	if err != nil {
		return errs.NewBadRequestError("invalid or expired reset token")
	}

	// Check if token is expired
	if u.PasswordResetExpiresAt != nil && u.PasswordResetExpiresAt.Before(time.Now()) {
		return errs.NewBadRequestError("reset token has expired")
	}

	// Hash new password
	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password and clear reset token
	err = s.userRepo.UpdatePassword(ctx, u.ID, passwordHash)
	if err != nil {
		return s.wrapError(err)
	}

	// Clear reset token
	err = s.userRepo.UpdatePasswordResetToken(ctx, u.ID, nil, nil)
	if err != nil {
		return s.wrapError(err)
	}

	return nil
}

// RefreshToken refreshes access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthResponse, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errs.NewUnauthorizedError("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errs.NewUnauthorizedError("invalid token type")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errs.NewBadRequestError("invalid user ID in token")
	}

	// Get user
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errs.NewUnauthorizedError("user not found")
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user.AuthResponse{
		User: &user.UserResponse{
			ID:            u.ID.String(),
			Name:          u.Name,
			Email:         u.Email,
			EmailVerified:   u.EmailVerified,
			OAuthProvider: u.OAuthProvider,
		},
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}

// TestAccountLogin logs in to test account
func (s *AuthService) TestAccountLogin(ctx context.Context) (*user.AuthResponse, error) {
	if s.testEmail == "" || s.testPassword == "" {
		return nil, errors.New("test account not configured")
	}

	// Find test user
	u, err := s.userRepo.GetByEmail(ctx, s.testEmail)
	if err != nil {
		return nil, errs.NewNotFoundError("test account not found")
	}

	// Verify password
	if u.PasswordHash == nil || !VerifyPassword(s.testPassword, *u.PasswordHash) {
		return nil, errs.NewUnauthorizedError("test account password mismatch")
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, u.ID)

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user.AuthResponse{
		User: &user.UserResponse{
			ID:            u.ID.String(),
			Name:          u.Name,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			OAuthProvider: u.OAuthProvider,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}

