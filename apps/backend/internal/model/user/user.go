package user

import (
	"time"

	"github.com/reche13/habitum/internal/model"
)

type User struct {
	model.Base

	Name                        string     `json:"name" db:"name"`
	Email                       string     `json:"email" db:"email"`
	PasswordHash                *string    `json:"-" db:"password_hash"`
	EmailVerified               bool       `json:"email_verified" db:"email_verified"`
	EmailVerificationToken      *string    `json:"-" db:"email_verification_token"`
	EmailVerificationExpiresAt  *time.Time `json:"-" db:"email_verification_expires_at"`
	PasswordResetToken          *string    `json:"-" db:"password_reset_token"`
	PasswordResetExpiresAt      *time.Time `json:"-" db:"password_reset_expires_at"`
	OAuthProvider               *string    `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthProviderID             *string    `json:"oauth_provider_id,omitempty" db:"oauth_provider_id"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}
