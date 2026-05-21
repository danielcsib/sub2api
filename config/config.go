// Package config provides configuration management for sub2api.
// It handles loading and validating application settings from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Host string
	Port int

	// Subscription settings
	SubURL      string
	RefreshInterval time.Duration
	UserAgent   string

	// API settings
	APIToken    string
	MaxConns    int

	// Cache settings
	CacheEnabled bool
	CacheTTL     time.Duration
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if any required configuration is missing or invalid.
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %w", err)
	}

	// Using 30 minutes as default refresh interval instead of 60 for more frequent updates
	refreshMinutes, err := strconv.Atoi(getEnv("REFRESH_INTERVAL", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_INTERVAL value: %w", err)
	}

	maxConns, err := strconv.Atoi(getEnv("MAX_CONNECTIONS", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_CONNECTIONS value: %w", err)
	}

	cacheTTLSeconds, err := strconv.Atoi(getEnv("CACHE_TTL", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL value: %w", err)
	}

	cacheEnabled, err := strconv.ParseBool(getEnv("CACHE_ENABLED", "true"))
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_ENABLED value: %w", err)
	}

	cfg := &Config{
		Host:            getEnv("HOST", "0.0.0.0"),
		Port:            port,
		SubURL:          getEnv("SUB_URL", ""),
		RefreshInterval: time.Duration(refreshMinutes) * time.Minute,
		UserAgent:       getEnv("USER_AGENT", "sub2api/1.0"),
		APIToken:        getEnv("API_TOKEN", ""),
		MaxConns:        maxConns,
		CacheEnabled:    cacheEnabled,
		CacheTTL:        time.Duration(cacheTTLSeconds) * time.Second,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required configuration values are present and valid.
func (c *Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535, got %d", c.Port)
	}
	if c.MaxConns < 1 {
		return fmt.Errorf("MAX_CONNECTIONS must be at least 1, got %d", c.MaxConns)
	}
	if c.RefreshInterval < time.Minute {
		return fmt.Errorf("REFRESH_INTERVAL must be at least 1 minute")
	}
	return nil
}

// Addr returns the full host:port address string for the server to listen on.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// getEnv retrieves an environment variable value, returning a default if not set.
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
