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
}

type ServerConfig struct {
	Port string `koanf:"port" validate:"required"`
}

type DatabaseConfig struct {
	URL string `koanf:"url" validate:"required"`
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
