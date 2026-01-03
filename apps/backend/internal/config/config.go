package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server   ServerConfig   `koanf:"server" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
	Auth     AuthConfig     `koanf:"auth" validate:"required"`
}

type ServerConfig struct {
	Port string `koanf:"port" validate:"required"`
}

type DatabaseConfig struct {
	URL string `koanf:"url" validate:"required"`
}

type AuthConfig struct {
	JWTSecret          string `koanf:"jwt_secret" validate:"required"`
	JWTAccessExpiry    string `koanf:"jwt_access_expiry"` // e.g., "15m"
	JWTRefreshExpiry   string `koanf:"jwt_refresh_expiry"` // e.g., "7d"
	ResendAPIKey       string `koanf:"resend_api_key"`
	GoogleClientID     string `koanf:"google_client_id"`
	GoogleClientSecret string `koanf:"google_client_secret"`
	FrontendURL        string `koanf:"frontend_url" validate:"required"`
	TestAccountEmail   string `koanf:"test_account_email"`
	TestAccountPassword string `koanf:"test_account_password"`
}


func Load() (*Config, error) {
	k := koanf.New(".")

	err := k.Load(env.Provider("HABITUM_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "HABITUM_"))
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	cfg := &Config{}

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal failed: %w", err)
	}

	validate := validator.New()
	
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
