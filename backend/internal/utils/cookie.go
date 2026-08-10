package utils

import (
	"music-player-backend/internal/logger"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CookieConfig struct {
	Secure   bool
	SameSite http.SameSite
	Domain   string
	HTTPOnly bool
	Path     string
}

func GetCookieConfig(env string) CookieConfig {
	domain := os.Getenv("COOKIE_DOMAIN")
	switch env {
	case "production":
		return CookieConfig{
			Secure:   true,
			SameSite: http.SameSiteNoneMode, // Untuk cross-site
			Domain:   domain,
			HTTPOnly: true,
			Path:     "/",
		}
	case "staging":
		return CookieConfig{
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Domain:   domain,
			HTTPOnly: true,
			Path:     "/",
		}
	default: // development
		return CookieConfig{
			Secure:   false, // Localhost tanpa HTTPS
			SameSite: http.SameSiteLaxMode,
			Domain:   "",
			HTTPOnly: true,
			Path:     "/",
		}
	}
}

func SetRefreshTokenCookie(c *gin.Context, refreshToken string, maxAge int) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	config := GetCookieConfig(env)

	c.SetSameSite(config.SameSite)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		config.Path,
		config.Domain,
		config.Secure,
		config.HTTPOnly,
	)

	// Untuk SameSite=None, Secure=true wajib
	if config.SameSite == http.SameSiteNoneMode && !config.Secure {
		logger.Log.Warn("Cookie SameSite=None requires Secure=true")
	}
}

func ClearRefreshTokenCookie(c *gin.Context) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	config := GetCookieConfig(env)

	c.SetSameSite(config.SameSite)
	c.SetCookie(
		"refresh_token",
		"",
		-1, // Expire immediately
		config.Path,
		config.Domain,
		config.Secure,
		config.HTTPOnly,
	)
}
