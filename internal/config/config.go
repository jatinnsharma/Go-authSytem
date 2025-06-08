
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment     string
	Port            string
	DatabaseURL     string
	RedisURL        string
	JWTSecret       string
	JWTExpiry       time.Duration
	RefreshExpiry   time.Duration
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	RateLimitPerMin int
	OTPExpiry       time.Duration
	BCryptCost      int
}

func Load() *Config {
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_EXPIRY", "168h"))
	otpExpiry, _ := time.ParseDuration(getEnv("OTP_EXPIRY", "5m"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT_PER_MIN", "100"))
	bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "12"))

	return &Config{
		Environment:     getEnv("ENVIRONMENT", "development"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost/authdb?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:       getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpiry:       jwtExpiry,
		RefreshExpiry:   refreshExpiry,
		SMTPHost:        getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:        smtpPort,
		SMTPUsername:    getEnv("SMTP_USERNAME", ""),
		SMTPPassword:    getEnv("SMTP_PASSWORD", ""),
		RateLimitPerMin: rateLimit,
		OTPExpiry:       otpExpiry,
		BCryptCost:      bcryptCost,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
