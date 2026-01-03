package service

import (
	"fmt"
	"time"

	"github.com/resendlabs/resend-go"
	"github.com/rs/zerolog"
)

type EmailService struct {
	client   *resend.Client
	from     string
	frontendURL string
	logger   zerolog.Logger
}

func NewEmailService(apiKey, from, frontendURL string, logger zerolog.Logger) *EmailService {
	client := resend.NewClient(apiKey)
	return &EmailService{
		client:      client,
		from:         from,
		frontendURL: frontendURL,
		logger:       logger,
	}
}

// SendVerificationEmail sends an email verification email
func (s *EmailService) SendVerificationEmail(email, name, token string) error {
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", s.frontendURL, token)
	
	subject := "Verify your Habitum account"
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Verify your email</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h1 style="color: #6366f1;">Welcome to Habitum!</h1>
				<p>Hi %s,</p>
				<p>Thank you for signing up! Please verify your email address by clicking the button below:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #6366f1; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">Verify Email</a>
				</div>
				<p>Or copy and paste this link into your browser:</p>
				<p style="word-break: break-all; color: #6366f1;">%s</p>
				<p>This link will expire in 24 hours.</p>
				<p>If you didn't create an account, you can safely ignore this email.</p>
			</div>
		</body>
		</html>
	`, name, verificationURL, verificationURL)

	plainBody := fmt.Sprintf(`
		Welcome to Habitum!
		
		Hi %s,
		
		Thank you for signing up! Please verify your email address by visiting:
		%s
		
		This link will expire in 24 hours.
		
		If you didn't create an account, you can safely ignore this email.
	`, name, verificationURL)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: subject,
		Html:    htmlBody,
		Text:    plainBody,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		s.logger.Error().Err(err).Str("email", email).Msg("failed to send verification email")
		return err
	}

	s.logger.Info().Str("email", email).Msg("verification email sent")
	return nil
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(email, name, token string) error {
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.frontendURL, token)
	
	subject := "Reset your Habitum password"
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Reset your password</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h1 style="color: #6366f1;">Reset Your Password</h1>
				<p>Hi %s,</p>
				<p>We received a request to reset your password. Click the button below to reset it:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #6366f1; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">Reset Password</a>
				</div>
				<p>Or copy and paste this link into your browser:</p>
				<p style="word-break: break-all; color: #6366f1;">%s</p>
				<p>This link will expire in 1 hour.</p>
				<p>If you didn't request a password reset, you can safely ignore this email.</p>
			</div>
		</body>
		</html>
	`, name, resetURL, resetURL)

	plainBody := fmt.Sprintf(`
		Reset Your Password
		
		Hi %s,
		
		We received a request to reset your password. Visit this link to reset it:
		%s
		
		This link will expire in 1 hour.
		
		If you didn't request a password reset, you can safely ignore this email.
	`, name, resetURL)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: subject,
		Html:    htmlBody,
		Text:    plainBody,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		s.logger.Error().Err(err).Str("email", email).Msg("failed to send password reset email")
		return err
	}

	s.logger.Info().Str("email", email).Msg("password reset email sent")
	return nil
}

// GetEmailVerificationExpiry returns the expiry time for email verification tokens (24 hours)
func GetEmailVerificationExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}

// GetPasswordResetExpiry returns the expiry time for password reset tokens (1 hour)
func GetPasswordResetExpiry() time.Time {
	return time.Now().Add(1 * time.Hour)
}

