package repositories

import (
	"context"

	"music-player-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepository interface {
	List(ctx context.Context) ([]models.Song, error)
	Search(ctx context.Context, query string) ([]models.Song, error)
	FindByID(ctx context.Context, id string) (*models.Song, error)
}

type songRepository struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) List(ctx context.Context) ([]models.Song, error) {
	query := `
		SELECT id, title, artist, album, cover_url, audio_url, duration, created_at
		FROM songs
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
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

func (r *songRepository) Search(ctx context.Context, query string) ([]models.Song, error) {
	sqlQuery := `
		SELECT id, title, artist, album, cover_url, audio_url, duration, created_at
		FROM songs
		WHERE title ILIKE $1 OR artist ILIKE $1
		ORDER BY created_at DESC
	`

	pattern := "%" + query + "%"
	rows, err := r.db.Query(ctx, sqlQuery, pattern)
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

func (r *songRepository) FindByID(ctx context.Context, id string) (*models.Song, error) {
	query := `
		SELECT id, title, artist, album, cover_url, audio_url, duration, created_at
		FROM songs
		WHERE id = $1
	`

	var s models.Song
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Title, &s.Artist, &s.Album, &s.CoverURL, &s.AudioURL, &s.Duration, &s.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}