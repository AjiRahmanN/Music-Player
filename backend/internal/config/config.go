package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
    DB DBConfig
    JWT JWTConfig
    Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	MaxConns int32
	MinConns int32
}


type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type ServerConfig struct {
	Env      string
	Port     string
	TimeZone string
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
    // Load and validate all configs
	// Load environment variables
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Parse durations
	// accessExpiry, _ := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRY"))
	accessExpiry, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY: %v", err)
	}

	// refreshExpirySTR, _ := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRY"))
	refreshExpiry, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY: %v", err)
	}
	config := &Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			MaxConns: 20,
			MinConns: 5,
		},
		JWT: JWTConfig{
			Secret:        os.Getenv("JWT_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Server: ServerConfig{
			Env:      env,
			Port:     os.Getenv("PORT"),
			TimeZone: "UTC",
		},
	}

	// Validate required configs
	if config.JWT.Secret == "" || config.JWT.Secret == "your_super_secret_key_change_in_production" {
		return nil, fmt.Errorf("JWT_SECRET must be set and secure")
	}

	if config.DB.Host == "" || config.DB.User == "" || config.DB.Name == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	if config.Server.Port == "" {
	return nil, fmt.Errorf("PORT must be set")
}

	AppConfig = config
	return config, nil
}