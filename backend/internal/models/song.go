package models

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

// Prefix path file hasil upload yang disajikan backend (lihat router /uploads).
const uploadsPrefix = "/uploads/"

var (
	assetMu      sync.RWMutex
	assetBaseURL string
)

// SetAssetBaseURL mengatur base URL publik backend (mis. https://api.contoh.com)
// yang dipakai untuk mengubah path upload relatif menjadi URL absolut saat
// dikirim ke client.
func SetAssetBaseURL(base string) {
	assetMu.Lock()
	defer assetMu.Unlock()
	assetBaseURL = strings.TrimRight(base, "/")
}

func getAssetBaseURL() string {
	assetMu.RLock()
	defer assetMu.RUnlock()
	return assetBaseURL
}

// ResolveAssetURL: "/uploads/audios/x.mp3" -> "<base>/uploads/audios/x.mp3".
// URL absolut (data lama atau link eksternal) dikembalikan apa adanya.
func ResolveAssetURL(u string) string {
	if strings.HasPrefix(u, uploadsPrefix) {
		if base := getAssetBaseURL(); base != "" {
			return base + u
		}
	}
	return u
}

// NormalizeAssetURL: kebalikan ResolveAssetURL. Kalau URL menunjuk ke file
// upload milik backend ini, simpan hanya path relatifnya sehingga data tidak
// rusak saat domain berubah. URL eksternal tidak diubah.
func NormalizeAssetURL(u string) string {
	u = strings.TrimSpace(u)
	if base := getAssetBaseURL(); base != "" && strings.HasPrefix(u, base+uploadsPrefix) {
		return strings.TrimPrefix(u, base)
	}
	return u
}

// MarshalJSON memastikan semua response yang memuat Song otomatis membawa URL
// absolut, tanpa perlu mengubah tiap handler.
func (s Song) MarshalJSON() ([]byte, error) {
	type alias Song
	a := alias(s)
	a.CoverURL = ResolveAssetURL(a.CoverURL)
	a.AudioURL = ResolveAssetURL(a.AudioURL)
	return json.Marshal(a)
}

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
