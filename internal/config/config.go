package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Port           string
	DBPath         string
	AuthToken      string
	BaseURL        string
	RetentionDays  int
	MaxBodySize    int64
	RateLimit      int
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:          getEnv("BACKFEEDR_PORT", "8080"),
		DBPath:        getEnv("BACKFEEDR_DB_PATH", "/data/backfeedr.db"),
		AuthToken:     getEnv("BACKFEEDR_AUTH_TOKEN", generateToken()),
		BaseURL:       getEnv("BACKFEEDR_BASE_URL", "http://localhost:8080"),
		RetentionDays: getEnvInt("BACKFEEDR_RETENTION_DAYS", 90),
		MaxBodySize:   getEnvInt64("BACKFEEDR_MAX_BODY_SIZE", 1*1024*1024), // 1MB
		RateLimit:     getEnvInt("BACKFEEDR_RATE_LIMIT", 100),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

func generateToken() string {
	// In production, this should generate a secure random token
	return "bf_admin_change_me_in_production"
}
