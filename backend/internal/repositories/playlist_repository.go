package repositories

import (
	"context"

	"music-player-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *models.Playlist) error
	ListByUser(ctx context.Context, userID string) ([]models.Playlist, error)
	ListAll(ctx context.Context) ([]models.Playlist, error)
	FindByID(ctx context.Context, id string) (*models.Playlist, error)
	Delete(ctx context.Context, id string) error
	AddSong(ctx context.Context, playlistID, songID string) error
	RemoveSong(ctx context.Context, playlistID, songID string) error
	GetSongs(ctx context.Context, playlistID string) ([]models.Song, error)
}

type playlistRepository struct {
	db *pgxpool.Pool
}

func NewPlaylistRepository(db *pgxpool.Pool) PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) Create(ctx context.Context, playlist *models.Playlist) error {
	query := `
		INSERT INTO playlists (user_id, name, description, cover_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, playlist.UserID, playlist.Name, playlist.Description, playlist.CoverURL).
		Scan(&playlist.ID, &playlist.CreatedAt)
}

func (r *playlistRepository) ListByUser(ctx context.Context, userID string) ([]models.Playlist, error) {
	query := `
		SELECT p.id, p.user_id, p.name, p.description, p.cover_url, p.created_at,
		       COUNT(ps.song_id) AS song_count
		FROM playlists p
		LEFT JOIN playlist_songs ps ON ps.playlist_id = p.id
		WHERE p.user_id = $1
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playlists := make([]models.Playlist, 0)
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CoverURL, &p.CreatedAt, &p.SongCount); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (r *playlistRepository) ListAll(ctx context.Context) ([]models.Playlist, error) {
	query := `
		SELECT p.id, p.user_id, p.name, p.description, p.cover_url, p.created_at,
		       COUNT(ps.song_id) AS song_count, u.username
		FROM playlists p
		LEFT JOIN playlist_songs ps ON ps.playlist_id = p.id
		JOIN users u ON u.id = p.user_id
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playlists := make([]models.Playlist, 0)
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CoverURL, &p.CreatedAt, &p.SongCount, &p.OwnerUsername); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (r *playlistRepository) FindByID(ctx context.Context, id string) (*models.Playlist, error) {
	query := `
		SELECT id, user_id, name, description, cover_url, created_at
		FROM playlists
		WHERE id = $1
	`

	var p models.Playlist
	err := r.db.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CoverURL, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *playlistRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM playlists WHERE id = $1`, id)
	return err
}

func (r *playlistRepository) AddSong(ctx context.Context, playlistID, songID string) error {
	query := `
		INSERT INTO playlist_songs (playlist_id, song_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, playlistID, songID)
	return err
}

func (r *playlistRepository) RemoveSong(ctx context.Context, playlistID, songID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2`, playlistID, songID)
	return err
}

func (r *playlistRepository) GetSongs(ctx context.Context, playlistID string) ([]models.Song, error) {
	query := `
		SELECT s.id, s.title, s.artist, s.album, s.cover_url, s.audio_url, s.duration, s.created_at
		FROM songs s
		JOIN playlist_songs ps ON ps.song_id = s.id
		WHERE ps.playlist_id = $1
		ORDER BY ps.added_at ASC
	`

	rows, err := r.db.Query(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := make([]models.Song, 0)
	for rows.Next() {
		var s models.Song
		if err := rows.Scan(&s.ID, &s.Title, &s.Artist, &s.Album, &s.CoverURL, &s.AudioURL, &s.Duration, &s.CreatedAt); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}