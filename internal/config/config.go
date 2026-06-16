// Package config loads runtime configuration from environment variables.
//
// SAFE: this file reads secrets from the environment, never hardcoded.
// It is the pattern the rest of the demo should follow.
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	StripeKey   string
	SlackToken  string
	OpenAIKey   string
	Debug       bool
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getenv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}
	return &Config{
		Port:        port,
		DatabaseURL: getenv("DATABASE_URL", ""),
		JWTSecret:   getenv("JWT_SECRET", ""),
		StripeKey:   getenv("STRIPE_SECRET_KEY", ""),
		SlackToken:  getenv("SLACK_BOT_TOKEN", ""),
		OpenAIKey:   getenv("OPENAI_API_KEY", ""),
		Debug:       getenv("DEBUG", "false") == "true",
	}, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
