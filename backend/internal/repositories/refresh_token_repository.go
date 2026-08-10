package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID, tokenHash string, expires time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (string, time.Time, error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteByUser(ctx context.Context, userID string) error
}

type refreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepo{db}
}

func (r *refreshTokenRepo) Create(
	ctx context.Context,
	userID, tokenHash string,
	expires time.Time,
) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		 VALUES ($1,$2,$3,$4)`,
		uuid.NewString(), userID, tokenHash, expires,
	)
	return err
}

func (r *refreshTokenRepo) FindByHash(
	ctx context.Context,
	tokenHash string,
) (string, time.Time, error) {

	var userID string
	var expires time.Time

	err := r.db.QueryRow(ctx,
		`SELECT user_id, expires_at
		 FROM refresh_tokens
		 WHERE token_hash=$1`,
		tokenHash,
	).Scan(&userID, &expires)

	if err != nil {
		return "", time.Time{}, err
	}

	return userID, expires, nil
}

func (r *refreshTokenRepo) DeleteByHash(
	ctx context.Context,
	tokenHash string,
) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE token_hash=$1`,
		tokenHash,
	)
	return err
}

func (r *refreshTokenRepo) DeleteByUser(
	ctx context.Context,
	userID string,
) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE user_id=$1`,
		userID,
	)
	return err
}
