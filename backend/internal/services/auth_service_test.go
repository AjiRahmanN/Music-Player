package services

import (
	"context"
	stdErrors "errors"
	"testing"
	"time"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/utils"
)

func newAuthFixture(t *testing.T) (*AuthService, *fakeRTRepo, *models.User) {
	t.Helper()
	utils.SetJWTSecret("test-secret-yang-cukup-panjang-untuk-hs256")
	hash, err := utils.HashPassword("Abcdef12")
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{ID: "u1", Username: "aji", Email: "aji@contoh.com", PasswordHash: hash, CreatedAt: time.Now()}
	rt := newFakeRTRepo()
	return NewAuthService(newFakeUserRepo(u), rt), rt, u
}

func TestValidateRegisterRequest(t *testing.T) {
	ok := models.RegisterRequest{Username: "aji_01", Email: "aji@contoh.com", Password: "Abcdef12"}
	if err := validateRegisterRequest(&ok); err != nil {
		t.Fatalf("request valid ditolak: %v", err)
	}
	cases := map[string]func(*models.RegisterRequest){
		"username kosong": func(r *models.RegisterRequest) { r.Username = "" },
		"username pendek": func(r *models.RegisterRequest) { r.Username = "ab" },
		"username simbol": func(r *models.RegisterRequest) { r.Username = "aji!" },
		"email kosong":    func(r *models.RegisterRequest) { r.Email = "" },
		"email invalid":   func(r *models.RegisterRequest) { r.Email = "bukan-email" },
		"password kosong": func(r *models.RegisterRequest) { r.Password = "" },
		"password lemah":  func(r *models.RegisterRequest) { r.Password = "abc" },
		"username terlalu panjang": func(r *models.RegisterRequest) {
			r.Username = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		},
	}
	for name, mut := range cases {
		r := ok
		mut(&r)
		if err := validateRegisterRequest(&r); err == nil {
			t.Errorf("%s: seharusnya error", name)
		}
	}
}

func TestLoginSuccessAndFailures(t *testing.T) {
	svc, rt, _ := newAuthFixture(t)
	ctx := context.Background()

	access, refresh, err := svc.Login(ctx, &models.LoginRequest{Email: " AJI@contoh.com ", Password: "Abcdef12"})
	if err != nil || access == "" || refresh == "" {
		t.Fatalf("login valid gagal: %v", err)
	}
	if rt.count() != 1 {
		t.Fatalf("refresh token tersimpan = %d, want 1", rt.count())
	}

	for name, req := range map[string]*models.LoginRequest{
		"password salah":      {Email: "aji@contoh.com", Password: "Salah123"},
		"email tak terdaftar": {Email: "nobody@contoh.com", Password: "Abcdef12"},
	} {
		_, _, err := svc.Login(ctx, req)
		if !stdErrors.Is(err, apperrors.ErrInvalidCredentials) {
			t.Errorf("%s: err = %v, want ErrInvalidCredentials (pesan harus seragam)", name, err)
		}
	}

	if _, _, err := svc.Login(ctx, &models.LoginRequest{}); err == nil {
		t.Error("kredensial kosong seharusnya error validasi")
	}
}

func TestRefreshRotatesAndOldTokenIsRejected(t *testing.T) {
	svc, rt, _ := newAuthFixture(t)
	ctx := context.Background()
	_, refresh, err := svc.Login(ctx, &models.LoginRequest{Email: "aji@contoh.com", Password: "Abcdef12"})
	if err != nil {
		t.Fatal(err)
	}

	newAccess, newRefresh, err := svc.Refresh(ctx, refresh)
	if err != nil || newAccess == "" || newRefresh == "" || newRefresh == refresh {
		t.Fatalf("refresh gagal / tidak dirotasi: %v", err)
	}
	if rt.count() != 1 {
		t.Fatalf("setelah rotasi harus tersisa 1 token, got %d", rt.count())
	}
	if _, _, err := svc.Refresh(ctx, refresh); err == nil {
		t.Fatal("refresh token lama seharusnya ditolak setelah rotasi")
	}
	if _, _, err := svc.Refresh(ctx, newRefresh); err != nil {
		t.Fatalf("token baru seharusnya valid: %v", err)
	}
}

func TestRefreshRejectsEmptyUnknownAndExpired(t *testing.T) {
	svc, rt, u := newAuthFixture(t)
	ctx := context.Background()

	if _, _, err := svc.Refresh(ctx, ""); err == nil {
		t.Error("token kosong seharusnya ditolak")
	}
	if _, _, err := svc.Refresh(ctx, "token-ngawur"); err == nil {
		t.Error("token tak dikenal seharusnya ditolak")
	}

	_ = rt.Create(ctx, u.ID, utils.HashToken("kedaluwarsa"), time.Now().Add(-time.Minute))
	if _, _, err := svc.Refresh(ctx, "kedaluwarsa"); !stdErrors.Is(err, apperrors.ErrTokenExpired) {
		t.Errorf("token expired: err = %v, want ErrTokenExpired", err)
	}
	if rt.count() != 0 {
		t.Error("token expired seharusnya dibersihkan dari DB")
	}
}

func TestLogoutOnlyRemovesCurrentSession(t *testing.T) {
	svc, rt, _ := newAuthFixture(t)
	ctx := context.Background()
	creds := &models.LoginRequest{Email: "aji@contoh.com", Password: "Abcdef12"}
	_, laptop, _ := svc.Login(ctx, creds)
	_, phone, _ := svc.Login(ctx, creds)
	if rt.count() != 2 {
		t.Fatalf("harus ada 2 sesi, got %d", rt.count())
	}

	if err := svc.Logout(ctx, laptop); err != nil {
		t.Fatal(err)
	}
	if rt.count() != 1 {
		t.Fatalf("setelah logout satu sesi harus tersisa 1, got %d", rt.count())
	}
	if _, _, err := svc.Refresh(ctx, laptop); err == nil {
		t.Error("sesi yang sudah logout tidak boleh bisa refresh")
	}
	if _, _, err := svc.Refresh(ctx, phone); err != nil {
		t.Errorf("sesi lain harus tetap hidup: %v", err)
	}
	if err := svc.Logout(ctx, ""); err != nil {
		t.Errorf("logout tanpa token tidak boleh error: %v", err)
	}
}

func TestLogoutAllRemovesEverySession(t *testing.T) {
	svc, rt, _ := newAuthFixture(t)
	ctx := context.Background()
	creds := &models.LoginRequest{Email: "aji@contoh.com", Password: "Abcdef12"}
	_, a, _ := svc.Login(ctx, creds)
	_, _, _ = svc.Login(ctx, creds)
	if err := svc.LogoutAll(ctx, a); err != nil {
		t.Fatal(err)
	}
	if rt.count() != 0 {
		t.Fatalf("LogoutAll harus menghapus semua sesi, sisa %d", rt.count())
	}
}

func TestRefreshTokenExpiryFollowsConfig(t *testing.T) {
	svc, rt, _ := newAuthFixture(t)
	ctx := context.Background()
	oldA, oldR := utils.AccessExpiry(), utils.RefreshExpiry()
	t.Cleanup(func() { utils.SetExpiry(oldA, oldR) })
	utils.SetExpiry(0, 2*time.Hour)

	_, refresh, err := svc.Login(ctx, &models.LoginRequest{Email: "aji@contoh.com", Password: "Abcdef12"})
	if err != nil {
		t.Fatal(err)
	}
	_, exp, _ := rt.FindByHash(ctx, utils.HashToken(refresh))
	if d := time.Until(exp); d > 2*time.Hour+time.Minute || d < 2*time.Hour-time.Minute {
		t.Fatalf("masa berlaku refresh token = %v, want ~2 jam", d)
	}
}
