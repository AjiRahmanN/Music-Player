package services

import (
	"context"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/repositories"
)

type SongService struct {
	repo repositories.SongRepository
}

func NewSongService(repo repositories.SongRepository) *SongService {
	return &SongService{repo}
}

// Batas pagination. Tanpa parameter, client mendapat DefaultSongLimit lagu.
const (
	DefaultSongLimit = 100
	MaxSongLimit     = 200
)

// NormalizePage merapikan limit/offset dari client ke rentang yang aman.
func NormalizePage(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = DefaultSongLimit
	}
	if limit > MaxSongLimit {
		limit = MaxSongLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (s *SongService) List(ctx context.Context, limit, offset int) ([]models.Song, error) {
	limit, offset = NormalizePage(limit, offset)
	songs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return songs, nil
}

func (s *SongService) Search(ctx context.Context, query string, limit, offset int) ([]models.Song, error) {
	limit, offset = NormalizePage(limit, offset)
	songs, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return songs, nil
}

func (s *SongService) Create(ctx context.Context, req *models.CreateSongRequest) (*models.Song, error) {
	if err := req.Validate(); err != nil {
		return nil, apperrors.NewValidationError("song", err.Error())
	}

	song := &models.Song{
		Title:    req.Title,
		Artist:   req.Artist,
		CoverURL: models.NormalizeAssetURL(req.CoverURL),
		AudioURL: models.NormalizeAssetURL(req.AudioURL),
		Duration: req.Duration,
	}
	if req.Album != "" {
		song.Album = &req.Album
	}

	if err := s.repo.Create(ctx, song); err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return song, nil
}

func (s *SongService) Update(ctx context.Context, id string, req *models.UpdateSongRequest) (*models.Song, error) {
	if err := req.Validate(); err != nil {
		return nil, apperrors.NewValidationError("song", err.Error())
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if existing == nil {
		return nil, apperrors.NewNotFoundError("song", id)
	}

	existing.Title = req.Title
	existing.Artist = req.Artist
	existing.CoverURL = models.NormalizeAssetURL(req.CoverURL)
	existing.AudioURL = models.NormalizeAssetURL(req.AudioURL)
	existing.Duration = req.Duration
	if req.Album != "" {
		existing.Album = &req.Album
	} else {
		existing.Album = nil
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return existing, nil
}

func (s *SongService) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if existing == nil {
		return apperrors.NewNotFoundError("song", id)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *SongService) GetByID(ctx context.Context, id string) (*models.Song, error) {
	song, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if song == nil {
		return nil, apperrors.NewNotFoundError("song", id)
	}
	return song, nil
}
