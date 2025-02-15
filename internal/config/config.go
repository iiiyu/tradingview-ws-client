package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
		DeviceToken string
		SessionID   string
		SessionSign string
	}
	Port string
}

func Load() (*Config, error) {
	// Initialize flags
	pflag.String("port", "3333", "Server port")
	pflag.String("db-host", "", "Database host")
	pflag.String("db-port", "", "Database port")
	pflag.String("db-user", "", "Database user")
	pflag.String("db-password", "", "Database password")
	pflag.String("db-name", "", "Database name")
	pflag.String("db-sslmode", "", "Database SSL mode")
	pflag.Parse()

	// Initialize Viper
	viper.SetEnvPrefix("") // This allows us to use environment variables without a prefix
	viper.AutomaticEnv()   // Automatically read environment variables

	// Bind flags to viper
	viper.BindPFlags(pflag.CommandLine)

	// Set up environment variable mappings
	viper.SetEnvPrefix("")
	viper.BindEnv("port", "PORT")
	viper.BindEnv("db-host", "DB_HOST")
	viper.BindEnv("db-port", "DB_PORT")
	viper.BindEnv("db-user", "DB_USER")
	viper.BindEnv("db-password", "DB_PASSWORD")
	viper.BindEnv("db-name", "DB_NAME")
	viper.BindEnv("db-sslmode", "DB_SSLMODE")
	viper.BindEnv("tradingview-device-token", "TRADINGVIEW_DEVICE_TOKEN")
	viper.BindEnv("tradingview-session-id", "TRADINGVIEW_SESSION_ID")
	viper.BindEnv("tradingview-session-sign", "TRADINGVIEW_SESSION_SIGN")

	cfg := &Config{}

	// Load server configuration
	cfg.Port = viper.GetString("port")

	// Load database configuration
	cfg.Database.Host = viper.GetString("db-host")
	cfg.Database.Port = viper.GetString("db-port")
	cfg.Database.User = viper.GetString("db-user")
	cfg.Database.Password = viper.GetString("db-password")
	cfg.Database.DBName = viper.GetString("db-name")
	cfg.Database.SSLMode = viper.GetString("db-sslmode")

	// Load TradingView configuration
	cfg.TradingView.DeviceToken = viper.GetString("tradingview-device-token")
	cfg.TradingView.SessionID = viper.GetString("tradingview-session-id")
	cfg.TradingView.SessionSign = viper.GetString("tradingview-session-sign")

	if cfg.TradingView.DeviceToken == "" {
		return nil, fmt.Errorf("TRADINGVIEW_DEVICE_TOKEN environment variable not set")
	}

	if cfg.TradingView.SessionID == "" {
		return nil, fmt.Errorf("TRADINGVIEW_SESSION_ID environment variable not set")
	}

	if cfg.TradingView.SessionSign == "" {
		return nil, fmt.Errorf("TRADINGVIEW_SESSION_SIGN environment variable not set")
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

func (c *Config) GetTradingViewConfig() (string, string, string) {
	return c.TradingView.DeviceToken, c.TradingView.SessionID, c.TradingView.SessionSign
}
