package service

import (
	"time"

	"github.com/reche13/habitum/internal/config"
	"github.com/reche13/habitum/internal/repository"
	"github.com/reche13/habitum/internal/sqlerr"
	"github.com/rs/zerolog"
)

type BaseService struct {
	resourceName string
}

func (b *BaseService) wrapError(err error) error {
	return sqlerr.WrapError(err, b.resourceName)
}

type Services struct{
	User *UserService
	Habit *HabitService
	HabitLog *HabitLogService
	Analytics *AnalyticsService
	Calendar *CalendarService
	Dashboard *DashboardService
	Auth *AuthService
}

func NewServices(repos *repository.Repositories, cfg *config.Config, logger zerolog.Logger) *Services {
	habitLogService := NewHabitLogService(repos.HabitLog)
	
	// Parse JWT expiry durations
	accessExpiry := 15 * time.Minute
	if cfg.Auth.JWTAccessExpiry != "" {
		if d, err := time.ParseDuration(cfg.Auth.JWTAccessExpiry); err == nil {
			accessExpiry = d
		}
	}
	
	refreshExpiry := 7 * 24 * time.Hour
	if cfg.Auth.JWTRefreshExpiry != "" {
		if d, err := time.ParseDuration(cfg.Auth.JWTRefreshExpiry); err == nil {
			refreshExpiry = d
		}
	}
	
	// Create JWT service
	jwtService := NewJWTService(cfg.Auth.JWTSecret, accessExpiry, refreshExpiry)
	
	// Create email service (optional - only if API key is provided)
	var emailService *EmailService
	if cfg.Auth.ResendAPIKey != "" {
		emailService = NewEmailService(cfg.Auth.ResendAPIKey, "noreply@habitum.app", cfg.Auth.FrontendURL, logger)
	}
	
	// Create OAuth service (optional - only if credentials are provided)
	var oauthService *OAuthService
	if cfg.Auth.GoogleClientID != "" && cfg.Auth.GoogleClientSecret != "" {
		oauthService = NewOAuthService(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.FrontendURL+"/auth/google/callback")
	}
	
	// Create auth service
	authService := NewAuthService(
		repos.User,
		jwtService,
		emailService,
		oauthService,
		logger,
		cfg.Auth.TestAccountEmail,
		cfg.Auth.TestAccountPassword,
	)
	
	return &Services{
		User: NewUserService(repos.User),
		Habit: NewHabitService(repos.Habit, habitLogService),
		HabitLog: habitLogService,
		Analytics: NewAnalyticsService(repos.Habit, repos.HabitLog),
		Calendar: NewCalendarService(repos.Habit, repos.HabitLog),
		Dashboard: NewDashboardService(repos.Habit, repos.HabitLog),
		Auth: authService,
	}
}