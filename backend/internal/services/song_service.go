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

func (s *SongService) List(ctx context.Context) ([]models.Song, error) {
	songs, err := s.repo.List(ctx)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return songs, nil
}

func (s *SongService) Search(ctx context.Context, query string) ([]models.Song, error) {
	songs, err := s.repo.Search(ctx, query)
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
		CoverURL: req.CoverURL,
		AudioURL: req.AudioURL,
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
	existing.CoverURL = req.CoverURL
	existing.AudioURL = req.AudioURL
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