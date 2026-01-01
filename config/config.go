package config

import (
	"os"
	"strconv"

	"github.com/leventeberry/goapi/logger"
)

// Config holds all application configuration
type Config struct {
	JWT struct {
		Secret          string
		ExpirationDays  int
	}
	RateLimit struct {
		RequestsPerMinute int
		BurstSize         int
	}
}

var AppConfig *Config

// Load reads configuration from environment variables with defaults
func Load() *Config {
	cfg := &Config{}

	// JWT Configuration
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		logger.Log.Fatal().Msg("JWT_SECRET environment variable is required")
	}

	expirationDaysStr := os.Getenv("JWT_EXPIRATION_DAYS")
	if expirationDaysStr == "" {
		cfg.JWT.ExpirationDays = 60 // Default: 60 days
	} else {
		days, err := strconv.Atoi(expirationDaysStr)
		if err != nil || days < 1 {
			logger.Log.Warn().Str("value", expirationDaysStr).Msg("Invalid JWT_EXPIRATION_DAYS, using default 60")
			cfg.JWT.ExpirationDays = 60
		} else {
			cfg.JWT.ExpirationDays = days
		}
	}

	// Rate Limit Configuration
	requestsPerMinuteStr := os.Getenv("RATE_LIMIT_REQUESTS_PER_MINUTE")
	if requestsPerMinuteStr == "" {
		cfg.RateLimit.RequestsPerMinute = 60 // Default: 60 requests per minute
	} else {
		rpm, err := strconv.Atoi(requestsPerMinuteStr)
		if err != nil || rpm < 1 {
			logger.Log.Warn().Str("value", requestsPerMinuteStr).Msg("Invalid RATE_LIMIT_REQUESTS_PER_MINUTE, using default 60")
			cfg.RateLimit.RequestsPerMinute = 60
		} else {
			cfg.RateLimit.RequestsPerMinute = rpm
		}
	}

	burstSizeStr := os.Getenv("RATE_LIMIT_BURST_SIZE")
	if burstSizeStr == "" {
		cfg.RateLimit.BurstSize = 10 // Default: burst size of 10
	} else {
		burst, err := strconv.Atoi(burstSizeStr)
		if err != nil || burst < 1 {
			logger.Log.Warn().Str("value", burstSizeStr).Msg("Invalid RATE_LIMIT_BURST_SIZE, using default 10")
			cfg.RateLimit.BurstSize = 10
		} else {
			cfg.RateLimit.BurstSize = burst
		}
	}

	AppConfig = cfg
	return cfg
}

// Get returns the global configuration instance
func Get() *Config {
	if AppConfig == nil {
		logger.Log.Fatal().Msg("Configuration not loaded. Call config.Load() first.")
	}
	return AppConfig
}

