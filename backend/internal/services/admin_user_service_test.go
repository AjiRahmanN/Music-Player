package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

func newAdminFixture() (*AdminUserService, *fakeAdminRepo) {
	repo := &fakeAdminRepo{users: map[string]*models.AdminUser{
		"admin1": {ID: "admin1", Username: "admin_satu", Email: "a1@x.com", IsAdmin: true},
		"admin2": {ID: "admin2", Username: "admin_dua", Email: "a2@x.com", IsAdmin: true},
		"user1":  {ID: "user1", Username: "biasa", Email: "u1@x.com"},
	}}
	return NewAdminUserService(repo), repo
}

func req(username, email string, admin bool) *models.UpdateUserByAdminRequest {
	return &models.UpdateUserByAdminRequest{Username: username, Email: email, IsAdmin: admin}
}

func appCode(err error) string {
	var app *apperrors.AppError
	if errors.As(err, &app) {
		return app.Code
	}
	return ""
}

func TestAdminCannotRemoveOwnAdminAccess(t *testing.T) {
	svc, repo := newAdminFixture()
	_, err := svc.Update(context.Background(), "admin1", "admin1", req("admin_satu", "a1@x.com", false))
	if appCode(err) != "VALIDATION_ERROR" {
		t.Fatalf("err = %v, want VALIDATION_ERROR", err)
	}
	if !repo.users["admin1"].IsAdmin {
		t.Fatal("hak admin tidak boleh berubah")
	}
}

func TestCannotRemoveLastAdmin(t *testing.T) {
	svc, repo := newAdminFixture()
	delete(repo.users, "admin2") // tinggal satu admin
	// admin lain (non-admin lagi pun) mencoba mencabut admin terakhir
	_, err := svc.Update(context.Background(), "someone-else", "admin1", req("admin_satu", "a1@x.com", false))
	if appCode(err) != "VALIDATION_ERROR" {
		t.Fatalf("err = %v, want VALIDATION_ERROR (admin terakhir)", err)
	}
}

func TestAdminCanDemoteAnotherAdminWhenOthersRemain(t *testing.T) {
	svc, repo := newAdminFixture()
	if _, err := svc.Update(context.Background(), "admin1", "admin2", req("admin_dua", "a2@x.com", false)); err != nil {
		t.Fatal(err)
	}
	if repo.users["admin2"].IsAdmin {
		t.Fatal("admin2 seharusnya sudah bukan admin")
	}
}

func TestUpdateNormalizesInput(t *testing.T) {
	svc, repo := newAdminFixture()
	if _, err := svc.Update(context.Background(), "admin1", "user1", req("  baru_nama  ", "  BARU@X.com ", false)); err != nil {
		t.Fatal(err)
	}
	u := repo.users["user1"]
	if u.Username != "baru_nama" || u.Email != "baru@x.com" {
		t.Fatalf("input tidak dinormalisasi: %+v", u)
	}
}

func TestUpdateMapsUniqueViolationToDuplicate(t *testing.T) {
	svc, repo := newAdminFixture()
	for constraint, field := range map[string]string{
		"users_email_lower_key":    "email",
		"users_username_lower_key": "username",
	} {
		repo.updateErr = fmt.Errorf("exec: %w", &pgconn.PgError{Code: "23505", ConstraintName: constraint})
		_, err := svc.Update(context.Background(), "admin1", "user1", req("nama_baru", "baru@x.com", false))
		if appCode(err) != "DUPLICATE_ERROR" {
			t.Fatalf("constraint %s (%s): err = %v, want DUPLICATE_ERROR", constraint, field, err)
		}
	}
	repo.updateErr = errors.New("koneksi putus")
	if _, err := svc.Update(context.Background(), "admin1", "user1", req("nama_baru", "baru@x.com", false)); appCode(err) != "INTERNAL_ERROR" {
		t.Fatalf("error lain harus INTERNAL_ERROR, got %v", err)
	}
}

func TestUpdateUnknownUser(t *testing.T) {
	svc, _ := newAdminFixture()
	if _, err := svc.Update(context.Background(), "admin1", "ghost", req("nama_baru", "baru@x.com", false)); !isNotFound(err) {
		t.Fatalf("err = %v, want NOT_FOUND", err)
	}
}

func TestAdminCannotDeleteSelf(t *testing.T) {
	svc, repo := newAdminFixture()
	if err := svc.Delete(context.Background(), "admin1", "admin1"); appCode(err) != "VALIDATION_ERROR" {
		t.Fatalf("err = %v", err)
	}
	if _, ok := repo.users["admin1"]; !ok {
		t.Fatal("akun sendiri tidak boleh terhapus")
	}
	if err := svc.Delete(context.Background(), "admin1", "user1"); err != nil {
		t.Fatal(err)
	}
}
