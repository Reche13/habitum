package user

// LoginRequest represents email/password login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// SignupRequest represents email/password signup request
type SignupRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// GoogleAuthRequest represents Google OAuth token request
type GoogleAuthRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"` // seconds
}

// UserResponse represents user data in auth response
type UserResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	OAuthProvider *string `json:"oauth_provider,omitempty"`
}
