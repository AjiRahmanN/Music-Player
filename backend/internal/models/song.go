package models

import "time"

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
