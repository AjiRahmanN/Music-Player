package models

import (
	"errors"
	"strings"
	"time"
)

// AdminUser dipakai khusus di endpoint admin — beda dari User biasa karena
// ikut membawa is_admin, yang tidak diekspos di endpoint /me atau /profile.
type AdminUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateUserByAdminRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

func (r *UpdateUserByAdminRequest) Validate() error {
	username := strings.TrimSpace(r.Username)
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 || len(username) > 50 {
		return errors.New("username must be 3-50 characters")
	}
	for _, ch := range username {
		if !(('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || ch == '_') {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	email := strings.TrimSpace(r.Email)
	if email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}