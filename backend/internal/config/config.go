package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
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
	// PublicBaseURL dipakai untuk membangun URL absolut file upload
	// (cover/audio) yang dikirim ke client, mis. https://api.contoh.com
	PublicBaseURL string
}

const (
	DefaultAccessExpiry  = 15 * time.Minute
	DefaultRefreshExpiry = 30 * 24 * time.Hour
)

var AppConfig *Config

// parseDurationEnv membaca durasi dari env; kosong -> fallback; invalid/<=0 -> error.
func parseDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %v", key, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0", key)
	}
	return d, nil
}

func LoadConfig() (*Config, error) {
	// Load and validate all configs
	// Load environment variables
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Durasi token: opsional. Kalau env kosong, pakai default yang aman
	// (15 menit untuk access token, 30 hari untuk refresh token).
	accessExpiry, err := parseDurationEnv("JWT_ACCESS_EXPIRY", DefaultAccessExpiry)
	if err != nil {
		return nil, err
	}
	refreshExpiry, err := parseDurationEnv("JWT_REFRESH_EXPIRY", DefaultRefreshExpiry)
	if err != nil {
		return nil, err
	}

	// PORT adalah nama utama; SERVER_PORT diterima sebagai fallback (dipakai
	// di README lama). Kalau dua-duanya kosong, jatuh ke 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("SERVER_PORT")
	}
	if port == "" {
		port = "8080"
	}
	publicBaseURL := strings.TrimRight(os.Getenv("PUBLIC_BASE_URL"), "/")
	if publicBaseURL == "" {
		publicBaseURL = "http://localhost:" + port
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
			Env:           env,
			Port:          port,
			TimeZone:      "UTC",
			PublicBaseURL: publicBaseURL,
		},
	}

	// Validate required configs
	if config.JWT.Secret == "" || config.JWT.Secret == "your_super_secret_key_change_in_production" {
		return nil, fmt.Errorf("JWT_SECRET must be set and secure")
	}

	if config.DB.Host == "" || config.DB.User == "" || config.DB.Name == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	AppConfig = config
	return config, nil
}
