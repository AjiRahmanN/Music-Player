package services

import (
	"context"
	"errors"
	"strings"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo}
}

// ======================
// GET CURRENT USER
// ======================
func (s *UserService) GetMe(
	ctx context.Context,
	userID string,
) (*models.User, error) {

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ======================
// UPDATE PROFILE
// ======================
func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID string,
	req *models.UpdateProfileRequest,
) (*models.User, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, apperrors.NewValidationError("profile", err.Error())
	}

	// Normalize input
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Ambil data user saat ini supaya field lain (password hash, dsb) tidak hilang
	current, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		var notFound *repositories.NotFoundError
		if errors.As(err, &notFound) {
			return nil, apperrors.NewNotFoundError("user", userID)
		}
		return nil, apperrors.NewInternalError(err)
	}
	if current == nil {
		return nil, apperrors.NewNotFoundError("user", userID)
	}

	// Tidak ada perubahan, tidak perlu tulis ke DB
	if current.Username == req.Username && current.Email == req.Email {
		return current, nil
	}

	current.Username = req.Username
	current.Email = req.Email

	if err := s.repo.Update(ctx, current); err != nil {
		var dupErr *repositories.DuplicateError
		if errors.As(err, &dupErr) {
			return nil, apperrors.NewDuplicateError(dupErr.Field, dupErr.Value)
		}
		var valErr *repositories.ValidationError
		if errors.As(err, &valErr) {
			return nil, apperrors.NewValidationError(valErr.Field, valErr.Message)
		}
		return nil, apperrors.NewInternalError(err)
	}

	return current, nil
}
