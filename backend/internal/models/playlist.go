package models

import (
	"errors"
	"strings"
	"time"
)

type Playlist struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CoverURL    *string   `json:"coverUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	// SongCount hanya diisi saat list playlist (bukan detail), lewat query JOIN
	SongCount int `json:"songCount"`
	// OwnerUsername hanya diisi di endpoint admin (lihat playlist semua user)
	OwnerUsername string `json:"ownerUsername,omitempty"`
}

// PlaylistDetail dipakai saat GET /api/playlists/:id, ikut membawa daftar lagunya
type PlaylistDetail struct {
	Playlist
	Songs []Song `json:"songs"`
}

type CreatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreatePlaylistRequest) Validate() error {
	name := strings.TrimSpace(r.Name)
	if name == "" {
		return errors.New("playlist name is required")
	}
	if len(name) > 255 {
		return errors.New("playlist name must be at most 255 characters")
	}
	return nil
}

type AddSongToPlaylistRequest struct {
	SongID string `json:"songId"`
}

func (r *AddSongToPlaylistRequest) Validate() error {
	if strings.TrimSpace(r.SongID) == "" {
		return errors.New("songId is required")
	}
	return nil
}