// File: internal/models/user.go
package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type User struct {
    ID           string    `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

// Validate - domain validation rules
func (u *User) Validate() error {
    if u.ID == "" {
        return errors.New("user ID cannot be empty")
    }
    
    if len(u.Username) < 3 || len(u.Username) > 50 {
        return errors.New("username must be 3-50 characters")
    }
    
    // Username pattern validation
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !usernameRegex.MatchString(u.Username) {
        return errors.New("username can only contain letters, numbers, and underscores")
    }
    
    // Email validation
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(u.Email) {
        return errors.New("invalid email format")
    }
    
    if u.CreatedAt.IsZero() {
        return errors.New("created_at cannot be zero")
    }
    
    return nil
}

// BeforeSave - hook untuk consistency check sebelum save
func (u *User) BeforeSave() error {
    // Normalize data
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    u.Username = strings.TrimSpace(u.Username)
    
    return u.Validate()
}