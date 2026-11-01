package utils

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ValidatePassword is the single password policy used by authentication flows.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasNumber bool
	for _, ch := range password {
		switch {
		case 'A' <= ch && ch <= 'Z':
			hasUpper = true
		case 'a' <= ch && ch <= 'z':
			hasLower = true
		case '0' <= ch && ch <= '9':
			hasNumber = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}

var (
	dummyHashOnce sync.Once
	dummyHash     []byte
)

// BurnPasswordCheck menjalankan satu perbandingan bcrypt terhadap hash palsu
// dengan cost yang sama seperti hash asli. Hasilnya dibuang; tujuannya hanya
// menyamakan durasi respons login untuk email yang tidak ada.
func BurnPasswordCheck(password string) {
	dummyHashOnce.Do(func() {
		dummyHash, _ = bcrypt.GenerateFromPassword([]byte("melodify-dummy-password"), 12)
	})
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
}
