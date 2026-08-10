package models

import (
	"errors"
	"strings"

	"music-player-backend/internal/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileRequest digunakan untuk endpoint PUT /api/profile
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Validate validates register request
func (r *RegisterRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}

	if len(r.Username) < 3 || len(r.Username) > 50 {
		return errors.New("username must be 3-50 characters")
	}

	// Simple username validation without regex
	for _, ch := range r.Username {
		if !(('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || ch == '_') {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	// Basic email validation without regex
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return errors.New("invalid email format")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return utils.ValidatePassword(r.Password)
}

// Validate validates login request
func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// Validate validates update profile request
func (r *UpdateProfileRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}

	if len(r.Username) < 3 || len(r.Username) > 50 {
		return errors.New("username must be 3-50 characters")
	}

	for _, ch := range r.Username {
		if !(('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || ch == '_') {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}
