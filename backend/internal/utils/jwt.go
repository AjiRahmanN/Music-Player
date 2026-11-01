package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret []byte

// Durasi token. Default-nya sama dengan perilaku lama (15 menit / 30 hari),
// dan bisa di-override lewat SetExpiry() saat startup (dari config/env).
var (
	accessExpiry  = 15 * time.Minute
	refreshExpiry = 30 * 24 * time.Hour
)

// SetExpiry mengatur durasi access & refresh token. Nilai <= 0 diabaikan.
func SetExpiry(access, refresh time.Duration) {
	if access > 0 {
		accessExpiry = access
	}
	if refresh > 0 {
		refreshExpiry = refresh
	}
}

func AccessExpiry() time.Duration  { return accessExpiry }
func RefreshExpiry() time.Duration { return refreshExpiry }

func SetJWTSecret(secret string) {
	JwtSecret = []byte(secret)
}

func InitJWTSecret() {
	if len(JwtSecret) == 0 {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" || secret == "your_super_secret_key_change_in_production" {
			log.Fatal("JWT_SECRET must be set in environment and be secure")
		}
		JwtSecret = []byte(secret)
	}
}

func GetJWTSecret() []byte {
	if len(JwtSecret) == 0 {
		InitJWTSecret()
	}
	return JwtSecret
}

// GenerateAccessToken creates a new JWT access token
func GenerateAccessToken(userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user ID cannot be empty")
	}

	if len(JwtSecret) == 0 {
		InitJWTSecret()
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(accessExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

// ParseAccessToken parses and validates a JWT access token
func ParseAccessToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("token string is empty")
	}

	if len(JwtSecret) == 0 {
		InitJWTSecret()
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	if claims["type"] != "access" {
		return "", fmt.Errorf("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}
