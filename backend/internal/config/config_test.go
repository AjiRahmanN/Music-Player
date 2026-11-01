package config

import (
	"testing"
	"time"
)

func setBaseEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_NAME", "n")
	t.Setenv("JWT_SECRET", "secret-panjang-dan-aman-0123456789")
	t.Setenv("JWT_ACCESS_EXPIRY", "")
	t.Setenv("JWT_REFRESH_EXPIRY", "")
	t.Setenv("PORT", "")
	t.Setenv("SERVER_PORT", "")
	t.Setenv("PUBLIC_BASE_URL", "")
}

func TestLoadConfigDefaultsWhenOptionalEnvMissing(t *testing.T) {
	setBaseEnv(t)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("env opsional kosong tidak boleh membuat startup gagal: %v", err)
	}
	if cfg.JWT.AccessExpiry != DefaultAccessExpiry || cfg.JWT.RefreshExpiry != DefaultRefreshExpiry {
		t.Fatalf("expiry default salah: %v / %v", cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	}
	if cfg.Server.Port != "8080" {
		t.Fatalf("port default = %q", cfg.Server.Port)
	}
	if cfg.Server.PublicBaseURL != "http://localhost:8080" {
		t.Fatalf("PublicBaseURL default = %q", cfg.Server.PublicBaseURL)
	}
}

func TestLoadConfigPortFallbackAndOverrides(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("SERVER_PORT", "9090") // nama lama di README
	t.Setenv("JWT_ACCESS_EXPIRY", "5m")
	t.Setenv("JWT_REFRESH_EXPIRY", "48h")
	t.Setenv("PUBLIC_BASE_URL", "https://api.contoh.com/")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("SERVER_PORT seharusnya jadi fallback, got %q", cfg.Server.Port)
	}
	if cfg.JWT.AccessExpiry != 5*time.Minute || cfg.JWT.RefreshExpiry != 48*time.Hour {
		t.Errorf("override expiry tidak diterapkan: %v / %v", cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	}
	if cfg.Server.PublicBaseURL != "https://api.contoh.com" {
		t.Errorf("PublicBaseURL = %q (trailing slash harus dibuang)", cfg.Server.PublicBaseURL)
	}

	t.Setenv("PORT", "7070") // PORT menang atas SERVER_PORT
	cfg, _ = LoadConfig()
	if cfg.Server.Port != "7070" {
		t.Errorf("PORT harus diprioritaskan, got %q", cfg.Server.Port)
	}
}

func TestLoadConfigRejectsInvalidValues(t *testing.T) {
	cases := map[string]map[string]string{
		"expiry bukan durasi": {"JWT_ACCESS_EXPIRY": "lima menit"},
		"expiry nol":          {"JWT_REFRESH_EXPIRY": "0s"},
		"expiry negatif":      {"JWT_ACCESS_EXPIRY": "-5m"},
		"secret kosong":       {"JWT_SECRET": ""},
		"secret placeholder":  {"JWT_SECRET": "your_super_secret_key_change_in_production"},
		"db tidak lengkap":    {"DB_HOST": ""},
	}
	for name, env := range cases {
		t.Run(name, func(t *testing.T) {
			setBaseEnv(t)
			for k, v := range env {
				t.Setenv(k, v)
			}
			if _, err := LoadConfig(); err == nil {
				t.Fatal("seharusnya error")
			}
		})
	}
}
