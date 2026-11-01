package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"music-player-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Custom error types untuk domain-specific errors
type (
	DuplicateError struct {
		Field string
		Value string
	}

	NotFoundError struct {
		Resource string
		ID       string
	}

	ValidationError struct {
		Field   string
		Message string
	}
)

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("%s '%s' already exists", e.Field, e.Value)
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// UserRepository interface dengan method yang lengkap
type UserRepository interface {
	// CRUD Operations
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error

	// Query Operations
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Count(ctx context.Context) (int64, error)

	// Bulk Operations
	FindAll(ctx context.Context, limit, offset int) ([]*models.User, error)
	FindByIDs(ctx context.Context, ids []string) ([]*models.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

// NewUserRepository membuat instance baru dengan dependency injection
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

// Create membuat user baru dengan validasi dan transaction
func (r *userRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// Validate input
	if err := validateUserForCreate(user); err != nil {
		return nil, err
	}

	// Normalize data
	normalizeUser(user)

	// Mulai transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Rollback tanpa syarat: aman dipanggil setelah Commit (hanya mengembalikan
	// pgx.ErrTxClosed yang kita abaikan). Sebelumnya rollback bergantung pada
	// variabel err, sehingga early-return (mis. DuplicateError) membocorkan
	// transaksi dan koneksi pool.
	defer func() { _ = tx.Rollback(ctx) }()

	// Cek duplicate dengan LOCK untuk menghindari race condition
	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE LOWER(email) = LOWER($1) 
			FOR UPDATE
		)`, user.Email).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate email: %w", err)
	}

	if exists {
		return nil, &DuplicateError{Field: "email", Value: user.Email}
	}

	// Cek duplicate username
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE LOWER(username) = LOWER($1)
		)`, user.Username).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate username: %w", err)
	}

	if exists {
		return nil, &DuplicateError{Field: "username", Value: user.Username}
	}

	// Set timestamps jika belum di-set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	// Insert user
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt)

	if err != nil {
		// Handle specific PostgreSQL errors
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505": // unique_violation
				// Extract detail dari error message
				if strings.Contains(pgErr.ConstraintName, "email") {
					return nil, &DuplicateError{Field: "email", Value: user.Email}
				} else if strings.Contains(pgErr.ConstraintName, "username") {
					return nil, &DuplicateError{Field: "username", Value: user.Username}
				}
			case "23514": // check_violation
				return nil, &ValidationError{Field: parseConstraintField(pgErr.ConstraintName), Message: pgErr.Message}
			}
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return created user (tanpa password hash untuk security)
	createdUser := &models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return createdUser, nil
}

// FindByID mencari user berdasarkan ID
func (r *userRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, &ValidationError{Field: "id", Message: "user ID cannot be empty"}
	}

	var user models.User

	// Query dengan context timeout handling
	err := r.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash, created_at
		FROM users 
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &NotFoundError{Resource: "user", ID: id}
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

// FindByEmail mencari user berdasarkan email (case-insensitive)
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, &ValidationError{Field: "email", Message: "email cannot be empty"}
	}

	var user models.User

	err := r.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash, created_at
		FROM users 
		WHERE LOWER(email) = LOWER($1)
	`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil tanpa error untuk login flow
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// ExistsByEmail mengecek apakah email sudah terdaftar
func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, &ValidationError{Field: "email", Message: "email cannot be empty"}
	}

	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE LOWER(email) = LOWER($1)
		)
	`, email).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// ExistsByUsername mengecek apakah username sudah terdaftar
func (r *userRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, &ValidationError{Field: "username", Message: "username cannot be empty"}
	}

	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE LOWER(username) = LOWER($1)
		)
	`, username).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}

// Update memperbarui data user
func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		return &ValidationError{Field: "id", Message: "user ID cannot be empty"}
	}

	// Normalize sebelum update
	normalizeUser(user)

	result, err := r.db.Exec(ctx, `
		UPDATE users 
		SET username = $1, email = $2, password_hash = $3
		WHERE id = $4
	`, user.Username, user.Email, user.PasswordHash, user.ID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return &DuplicateError{Field: "email", Value: user.Email}
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Cek jika ada row yang terupdate
	if result.RowsAffected() == 0 {
		return &NotFoundError{Resource: "user", ID: user.ID}
	}

	return nil
}

// Delete menghapus user berdasarkan ID
func (r *userRepo) Delete(ctx context.Context, id string) error {
	if id == "" {
		return &ValidationError{Field: "id", Message: "user ID cannot be empty"}
	}

	result, err := r.db.Exec(ctx, `
		DELETE FROM users 
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return &NotFoundError{Resource: "user", ID: id}
	}

	return nil
}

// Count menghitung total user
func (r *userRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// FindAll mengambil semua user dengan pagination
func (r *userRepo) FindAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit < 1 || limit > 100 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, username, email, created_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

// FindByIDs mencari multiple users berdasarkan IDs
func (r *userRepo) FindByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	if len(ids) == 0 {
		return []*models.User{}, nil
	}

	// Build query dengan IN clause
	query := `
		SELECT id, username, email, created_at
		FROM users 
		WHERE id = ANY($1)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by IDs: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

// ======================
// HELPER FUNCTIONS
// ======================

// validateUserForCreate validasi user sebelum create
func validateUserForCreate(user *models.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	if user.ID == "" {
		return &ValidationError{Field: "id", Message: "user ID cannot be empty"}
	}

	if user.Username == "" {
		return &ValidationError{Field: "username", Message: "username cannot be empty"}
	}

	if len(user.Username) < 3 || len(user.Username) > 50 {
		return &ValidationError{Field: "username", Message: "username must be between 3-50 characters"}
	}

	if user.Email == "" {
		return &ValidationError{Field: "email", Message: "email cannot be empty"}
	}

	if !strings.Contains(user.Email, "@") {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}

	if user.PasswordHash == "" {
		return &ValidationError{Field: "password", Message: "password hash cannot be empty"}
	}

	return nil
}

// normalizeUser normalisasi data user (lowercase email, trim spaces)
func normalizeUser(user *models.User) {
	if user == nil {
		return
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)
}

// parseConstraintField mengekstrak field name dari constraint error
func parseConstraintField(constraintName string) string {
	// Contoh: "users_email_key" -> "email"
	if strings.Contains(constraintName, "email") {
		return "email"
	}
	if strings.Contains(constraintName, "username") {
		return "username"
	}
	return "unknown"
}
