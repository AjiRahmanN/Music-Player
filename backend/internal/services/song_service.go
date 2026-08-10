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