package services

import (
	"context"
	stdErrors "errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/repositories"
)

type AdminUserService struct {
	repo repositories.AdminUserRepository
}

func NewAdminUserService(repo repositories.AdminUserRepository) *AdminUserService {
	return &AdminUserService{repo}
}

func (s *AdminUserService) List(ctx context.Context) ([]models.AdminUser, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return users, nil
}

// Update memperbarui user. requesterID adalah admin yang sedang melakukan
// perubahan; dipakai untuk mencegah admin mencabut hak admin dirinya sendiri
// dan mencegah sistem kehilangan admin terakhir.
func (s *AdminUserService) Update(ctx context.Context, requesterID, id string, req *models.UpdateUserByAdminRequest) (*models.AdminUser, error) {
	if err := req.Validate(); err != nil {
		return nil, apperrors.NewValidationError("user", err.Error())
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if existing == nil {
		return nil, apperrors.NewNotFoundError("user", id)
	}

	// Pencabutan hak admin: tidak boleh pada diri sendiri, dan tidak boleh
	// meninggalkan sistem tanpa admin.
	if existing.IsAdmin && !req.IsAdmin {
		if requesterID == id {
			return nil, apperrors.NewValidationError("user", "you cannot remove your own admin access")
		}
		count, err := s.repo.CountAdmins(ctx)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		if count <= 1 {
			return nil, apperrors.NewValidationError("user", "cannot remove the last admin")
		}
	}

	existing.Username = req.Username
	existing.Email = req.Email
	existing.IsAdmin = req.IsAdmin

	if err := s.repo.Update(ctx, existing); err != nil {
		if dup := duplicateFromPgError(err, req); dup != nil {
			return nil, dup
		}
		return nil, apperrors.NewInternalError(err)
	}

	return existing, nil
}

// Delete menghapus user. requesterID dicek supaya admin tidak bisa
// menghapus akunnya sendiri secara tidak sengaja lewat panel admin.
func (s *AdminUserService) Delete(ctx context.Context, requesterID, targetID string) error {
	if requesterID == targetID {
		return apperrors.NewValidationError("user", "you cannot delete your own account")
	}

	existing, err := s.repo.FindByID(ctx, targetID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if existing == nil {
		return apperrors.NewNotFoundError("user", targetID)
	}

	if err := s.repo.Delete(ctx, targetID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

// duplicateFromPgError menerjemahkan pelanggaran unique constraint PostgreSQL
// (SQLSTATE 23505) menjadi DuplicateError (HTTP 409) berdasarkan nama
// constraint/index (mengandung "email" atau "username").
func duplicateFromPgError(err error, req *models.UpdateUserByAdminRequest) *apperrors.AppError {
	var pgErr *pgconn.PgError
	if !stdErrors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}
	name := strings.ToLower(pgErr.ConstraintName)
	switch {
	case strings.Contains(name, "email"):
		return apperrors.NewDuplicateError("email", req.Email)
	case strings.Contains(name, "username"):
		return apperrors.NewDuplicateError("username", req.Username)
	}
	return nil
}
