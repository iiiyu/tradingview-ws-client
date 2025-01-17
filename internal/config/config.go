package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	TradingView struct {
		AuthToken string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Load database configuration
	cfg.Database.Host = getEnvOrDefault("DB_HOST", "192.168.1.48")
	cfg.Database.Port = getEnvOrDefault("DB_PORT", "6543")
	cfg.Database.User = getEnvOrDefault("DB_USER", "postgres")
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", "uUE1yOke9wIqSAwL7bZBfKJHb5WqDnzmPIc0tlg9rF86hb5m7djpKDHulKmGy3Iy")
	cfg.Database.DBName = getEnvOrDefault("DB_NAME", "postgres")
	cfg.Database.SSLMode = getEnvOrDefault("DB_SSLMODE", "disable")

	// Load TradingView configuration
	cfg.TradingView.AuthToken = os.Getenv("TRADINGVIEW_AUTH_TOKEN")
	if cfg.TradingView.AuthToken == "" {
		return nil, fmt.Errorf("TRADINGVIEW_AUTH_TOKEN environment variable not set")
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.DBName,
		c.Database.Password,
		c.Database.SSLMode,
	)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
