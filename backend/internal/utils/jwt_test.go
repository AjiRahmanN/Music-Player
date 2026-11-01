package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func setup(t *testing.T) {
	t.Helper()
	SetJWTSecret("test-secret-yang-cukup-panjang-untuk-hs256")
	oldA, oldR := accessExpiry, refreshExpiry
	t.Cleanup(func() { accessExpiry, refreshExpiry = oldA, oldR })
}

func TestAccessTokenRoundTrip(t *testing.T) {
	setup(t)
	tok, err := GenerateAccessToken("user-1")
	if err != nil {
		t.Fatal(err)
	}
	id, err := ParseAccessToken(tok)
	if err != nil || id != "user-1" {
		t.Fatalf("ParseAccessToken = %q, %v", id, err)
	}
}

func TestGenerateAccessTokenRejectsEmptyUser(t *testing.T) {
	setup(t)
	if _, err := GenerateAccessToken(""); err == nil {
		t.Fatal("user ID kosong seharusnya error")
	}
}

func TestParseAccessTokenRejectsWrongSecret(t *testing.T) {
	setup(t)
	tok, _ := GenerateAccessToken("user-1")
	SetJWTSecret("secret-lain-yang-juga-panjang-sekali-0000")
	if _, err := ParseAccessToken(tok); err == nil {
		t.Fatal("token dengan secret berbeda seharusnya ditolak")
	}
}

func TestParseAccessTokenRejectsWrongType(t *testing.T) {
	setup(t)
	claims := jwt.MapClaims{"user_id": "u", "type": "refresh", "exp": time.Now().Add(time.Hour).Unix()}
	tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtSecret)
	if _, err := ParseAccessToken(tok); err == nil {
		t.Fatal("token bertipe refresh seharusnya ditolak sebagai access token")
	}
}

func TestParseAccessTokenRejectsExpired(t *testing.T) {
	setup(t)
	SetExpiry(time.Millisecond, 0)
	tok, _ := GenerateAccessToken("user-1")
	time.Sleep(1100 * time.Millisecond) // exp dalam detik; tunggu lewat 1 detik
	if _, err := ParseAccessToken(tok); err == nil {
		t.Fatal("token kedaluwarsa seharusnya ditolak")
	}
}

func TestParseAccessTokenRejectsNoneAlg(t *testing.T) {
	setup(t)
	claims := jwt.MapClaims{"user_id": "u", "type": "access", "exp": time.Now().Add(time.Hour).Unix()}
	tok, _ := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if _, err := ParseAccessToken(tok); err == nil {
		t.Fatal("alg none seharusnya ditolak")
	}
}

func TestSetExpiryIgnoresNonPositive(t *testing.T) {
	setup(t)
	SetExpiry(10*time.Minute, 48*time.Hour)
	if AccessExpiry() != 10*time.Minute || RefreshExpiry() != 48*time.Hour {
		t.Fatal("SetExpiry tidak menerapkan nilai valid")
	}
	SetExpiry(0, -1)
	if AccessExpiry() != 10*time.Minute || RefreshExpiry() != 48*time.Hour {
		t.Fatal("SetExpiry seharusnya mengabaikan nilai <= 0")
	}
}

func TestHashTokenDeterministic(t *testing.T) {
	if HashToken("abc") != HashToken("abc") || HashToken("abc") == HashToken("abd") {
		t.Fatal("HashToken harus deterministik dan berbeda untuk input berbeda")
	}
}
