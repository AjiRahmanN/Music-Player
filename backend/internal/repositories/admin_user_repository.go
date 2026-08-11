package repositories

import (
	"context"

	"music-player-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminUserRepository interface {
	List(ctx context.Context) ([]models.AdminUser, error)
	FindByID(ctx context.Context, id string) (*models.AdminUser, error)
	Update(ctx context.Context, user *models.AdminUser) error
	Delete(ctx context.Context, id string) error
}

type adminUserRepository struct {
	db *pgxpool.Pool
}

func NewAdminUserRepository(db *pgxpool.Pool) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) List(ctx context.Context) ([]models.AdminUser, error) {
	query := `
		SELECT id, username, email, is_admin, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.AdminUser, 0)
	for rows.Next() {
		var u models.AdminUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.IsAdmin, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *adminUserRepository) FindByID(ctx context.Context, id string) (*models.AdminUser, error) {
	query := `SELECT id, username, email, is_admin, created_at FROM users WHERE id = $1`

	var u models.AdminUser
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email, &u.IsAdmin, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *adminUserRepository) Update(ctx context.Context, user *models.AdminUser) error {
	query := `UPDATE users SET username = $1, email = $2, is_admin = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.IsAdmin, user.ID)
	return err
}

func (r *adminUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}