package services

import (
	"context"

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

func (s *AdminUserService) Update(ctx context.Context, id string, req *models.UpdateUserByAdminRequest) (*models.AdminUser, error) {
	if err := req.Validate(); err != nil {
		return nil, apperrors.NewValidationError("user", err.Error())
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if existing == nil {
		return nil, apperrors.NewNotFoundError("user", id)
	}

	existing.Username = req.Username
	existing.Email = req.Email
	existing.IsAdmin = req.IsAdmin

	if err := s.repo.Update(ctx, existing); err != nil {
		// Catatan: kalau username/email bentrok dengan user lain, error mentah
		// dari database (unique constraint) akan lolos sebagai internal_error
		// generik, bukan pesan "sudah dipakai" yang rapi — karena tipe error
		// khusus duplikat di project ini terikat ke UserRepository yang tidak
		// disentuh oleh AdminUserRepository (sengaja, demi keamanan kode lama).
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