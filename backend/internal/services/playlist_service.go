package services

import (
	"context"
	"strings"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/repositories"
)

type PlaylistService struct {
	repo     repositories.PlaylistRepository
	songRepo repositories.SongRepository
}

func NewPlaylistService(repo repositories.PlaylistRepository, songRepo repositories.SongRepository) *PlaylistService {
	return &PlaylistService{repo: repo, songRepo: songRepo}
}

func (s *PlaylistService) Create(ctx context.Context, userID string, req *models.CreatePlaylistRequest) (*models.Playlist, error) {
	if err := req.Validate(); err != nil {
		return nil, apperrors.NewValidationError("playlist", err.Error())
	}

	name := strings.TrimSpace(req.Name)
	var description *string
	if trimmed := strings.TrimSpace(req.Description); trimmed != "" {
		description = &trimmed
	}

	playlist := &models.Playlist{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(ctx, playlist); err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	return playlist, nil
}

func (s *PlaylistService) ListByUser(ctx context.Context, userID string) ([]models.Playlist, error) {
	playlists, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return playlists, nil
}

// getOwnedPlaylist mengambil playlist dan sekaligus memastikan playlist itu
// benar milik userID. Kalau tidak ditemukan ATAU bukan miliknya, sengaja
// dikembalikan error yang sama (NotFound) supaya tidak membocorkan info
// bahwa playlist itu sebenarnya ada tapi milik orang lain.
func (s *PlaylistService) getOwnedPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	playlist, err := s.repo.FindByID(ctx, playlistID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if playlist == nil || playlist.UserID != userID {
		return nil, apperrors.NewNotFoundError("playlist", playlistID)
	}
	return playlist, nil
}

func (s *PlaylistService) GetDetail(ctx context.Context, userID, playlistID string) (*models.PlaylistDetail, error) {
	playlist, err := s.getOwnedPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	songs, err := s.repo.GetSongs(ctx, playlistID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	return &models.PlaylistDetail{
		Playlist: *playlist,
		Songs:    songs,
	}, nil
}

func (s *PlaylistService) Delete(ctx context.Context, userID, playlistID string) error {
	if _, err := s.getOwnedPlaylist(ctx, userID, playlistID); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, playlistID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *PlaylistService) AddSong(ctx context.Context, userID, playlistID string, req *models.AddSongToPlaylistRequest) error {
	if err := req.Validate(); err != nil {
		return apperrors.NewValidationError("playlist_song", err.Error())
	}
	if _, err := s.getOwnedPlaylist(ctx, userID, playlistID); err != nil {
		return err
	}

	song, err := s.songRepo.FindByID(ctx, req.SongID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if song == nil {
		return apperrors.NewNotFoundError("song", req.SongID)
	}

	if err := s.repo.AddSong(ctx, playlistID, req.SongID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *PlaylistService) RemoveSong(ctx context.Context, userID, playlistID, songID string) error {
	if _, err := s.getOwnedPlaylist(ctx, userID, playlistID); err != nil {
		return err
	}
	if err := s.repo.RemoveSong(ctx, playlistID, songID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}