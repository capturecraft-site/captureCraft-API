package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds runtime settings.
type Config struct {
	Port      int
	JWTSecret string
	TokenTTL  time.Duration
	Mode      string // server or lambda
}

// Load reads environment variables with defaults.
func Load() (Config, error) {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			port = parsed
		}
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	tokenTTL := 24 * time.Hour
	if v := os.Getenv("TOKEN_TTL_HOURS"); v != "" {
		hours, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid TOKEN_TTL_HOURS: %w", err)
		}
		tokenTTL = time.Duration(hours) * time.Hour
	}
	mode := os.Getenv("RUN_MODE")
	if mode == "" {
		mode = "server"
	}
	return Config{Port: port, JWTSecret: secret, TokenTTL: tokenTTL, Mode: mode}, nil
}
