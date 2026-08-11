package models

import (
	"errors"
	"time"
)

type Song struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Artist    string    `json:"artist"`
	Album     *string   `json:"album,omitempty"`
	CoverURL  string    `json:"coverUrl"`
	AudioURL  string    `json:"audioUrl"`
	Duration  int       `json:"duration"` // dalam detik
	CreatedAt time.Time `json:"createdAt"`
}

// CreateSongRequest dipakai handler admin untuk membuat/mengedit lagu.
type CreateSongRequest struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	CoverURL string `json:"coverUrl"`
	AudioURL string `json:"audioUrl"`
	Duration int    `json:"duration"`
}

func (r *CreateSongRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if r.Artist == "" {
		return errors.New("artist is required")
	}
	if r.AudioURL == "" {
		return errors.New("audioUrl is required")
	}
	if r.Duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	return nil
}

// UpdateSongRequest bentuknya sama persis dengan CreateSongRequest
type UpdateSongRequest = CreateSongRequest