package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/repositories"
	"music-player-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	// Import stdlib errors separately to avoid conflict

	stdErrors "errors"
)

type AuthService struct {
	userRepo repositories.UserRepository
	rtRepo   repositories.RefreshTokenRepository
}

func NewAuthService(
	userRepo repositories.UserRepository,
	rtRepo repositories.RefreshTokenRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		rtRepo:   rtRepo,
	}
}

// Register registers a new user with validation and returns the created user
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Validate input
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Normalize input
	normalizeRegisterRequest(req)

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("failed to check email existence: %w", err))
	}
	if exists {
		return nil, apperrors.NewDuplicateError("email", req.Email)
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("failed to check username existence: %w", err))
	}
	if exists {
		return nil, apperrors.NewDuplicateError("username", req.Username)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("failed to hash password: %w", err))
	}

	// Create user object
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now().UTC(),
	}

	// Save to repository
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		// Handle repository-specific errors
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "email") {
				return nil, apperrors.NewDuplicateError("email", req.Email)
			} else if strings.Contains(err.Error(), "username") {
				return nil, apperrors.NewDuplicateError("username", req.Username)
			}
		}
		return nil, apperrors.NewInternalError(fmt.Errorf("failed to create user: %w", err))
	}

	return createdUser, nil
}

// Login authenticates user and returns access and refresh tokens
func (s *AuthService) Login(
	ctx context.Context,
	req *models.LoginRequest,
) (accessToken string, refreshToken string, err error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return "", "", apperrors.NewValidationError("credentials", "email and password are required")
	}

	// Normalize email (lowercase)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Don't leak information about user existence
		return "", "", apperrors.ErrInvalidCredentials
	}

	// Check if user exists
	if user == nil {
		// Use same error message to prevent user enumeration
		return "", "", apperrors.ErrInvalidCredentials
	}

	// Verify password
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		// Use same error message regardless of password mismatch
		return "", "", apperrors.ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err = utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", apperrors.NewInternalError(fmt.Errorf("failed to generate access token: %w", err))
	}

	// Generate refresh token
	refreshToken, err = s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", apperrors.NewInternalError(fmt.Errorf("failed to generate refresh token: %w", err))
	}

	return accessToken, refreshToken, nil
}

// Refresh generates new access token using refresh token
func (s *AuthService) Refresh(
	ctx context.Context,
	rawRefreshToken string,
) (newAccessToken string, newRefreshToken string, err error) {
	// Validate input
	if rawRefreshToken == "" {
		return "", "", apperrors.NewAuthError("Refresh token is required")
	}

	// Hash the token to find in database
	tokenHash := utils.HashToken(rawRefreshToken)

	// Find token in database
	userID, expiresAt, err := s.rtRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return "", "", apperrors.NewAuthError("Invalid or expired refresh token")
		}
		return "", "", apperrors.NewInternalError(fmt.Errorf("failed to validate refresh token: %w", err))
	}

	// Check if token is expired
	if time.Now().After(expiresAt) {
		// Clean up expired token
		_ = s.rtRepo.DeleteByHash(ctx, tokenHash)
		return "", "", apperrors.ErrTokenExpired
	}

	// Token rotation: delete old refresh token
	if err := s.rtRepo.DeleteByHash(ctx, tokenHash); err != nil {
		// Log but continue (token might have been deleted by concurrent request)
		// Don't expose internal errors to client
	}

	// Generate new access token
	newAccessToken, err = utils.GenerateAccessToken(userID)
	if err != nil {
		return "", "", apperrors.NewInternalError(fmt.Errorf("failed to generate new access token: %w", err))
	}

	// Generate new refresh token
	newRefreshToken, err = s.generateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", apperrors.NewInternalError(fmt.Errorf("failed to generate new refresh token: %w", err))
	}

	return newAccessToken, newRefreshToken, nil
}

// LogoutAll invalidates all refresh tokens for a user
func (s *AuthService) LogoutAll(
	ctx context.Context,
	rawRefreshToken string,
) error {
	// If no token provided, nothing to do
	if rawRefreshToken == "" {
		return nil
	}

	// Hash the token
	tokenHash := utils.HashToken(rawRefreshToken)

	// Find user associated with this token
	userID, _, err := s.rtRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		// Token might already be invalid, which is fine for logout
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return apperrors.NewInternalError(fmt.Errorf("failed to validate refresh token: %w", err))
	}

	// Delete all refresh tokens for this user
	if err := s.rtRepo.DeleteByUser(ctx, userID); err != nil {
		return apperrors.NewInternalError(fmt.Errorf("failed to logout user: %w", err))
	}

	return nil
}

// GetUserFromToken validates token and returns user
func (s *AuthService) GetUserFromToken(
	ctx context.Context,
	accessToken string,
) (*models.User, error) {
	// Parse and validate JWT
	userID, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return nil, apperrors.ErrTokenInvalid
	}

	// Get user from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("failed to find user: %w", err))
	}

	if user == nil {
		return nil, apperrors.NewNotFoundError("user", "")
	}

	return user, nil
}

// ======================
// HELPER FUNCTIONS
// ======================

// validateRegisterRequest validates register request
func validateRegisterRequest(req *models.RegisterRequest) error {
	if req.Username == "" {
		return apperrors.NewValidationError("username", "is required")
	}

	if len(req.Username) < 3 {
		return apperrors.NewValidationError("username", "must be at least 3 characters")
	}

	if len(req.Username) > 50 {
		return apperrors.NewValidationError("username", "must be 50 characters or fewer")
	}

	// Username pattern (alphanumeric and underscore)
	if !usernameRegex.MatchString(req.Username) {
		return apperrors.NewValidationError("username", "can only contain letters, numbers, and underscores")
	}

	if req.Email == "" {
		return apperrors.NewValidationError("email", "is required")
	}

	// Basic email validation
	if !emailRegex.MatchString(req.Email) {
		return apperrors.NewValidationError("email", "invalid format")
	}

	if req.Password == "" {
		return apperrors.NewValidationError("password", "is required")
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		return apperrors.NewValidationError("password", err.Error())
	}
	return nil
}

// normalizeRegisterRequest normalizes register request data
func normalizeRegisterRequest(req *models.RegisterRequest) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	// Password trimming depends on security policy
	// Some passwords may intentionally have spaces
}

// generateRefreshToken creates a new refresh token for a user
func (s *AuthService) generateRefreshToken(
	ctx context.Context,
	userID string,
) (string, error) {
	// Generate raw refresh token
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", apperrors.NewInternalError(fmt.Errorf("failed to generate refresh token: %w", err))
	}
	rawRefreshToken := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Hash the token for storage
	tokenHash := utils.HashToken(rawRefreshToken)

	// Set expiration (30 days)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	// Store in database
	if err := s.rtRepo.Create(ctx, userID, tokenHash, expiresAt); err != nil {
		return "", apperrors.NewInternalError(fmt.Errorf("failed to store refresh token: %w", err))
	}

	return rawRefreshToken, nil
}

// isDuplicateError checks if error is a duplicate error
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	return strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "23505") // PostgreSQL unique violation code
}

// Regular expressions for validation (without lookahead)
var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)
