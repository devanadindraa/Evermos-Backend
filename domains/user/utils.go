package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

func comparePassword(storedHash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plain))
	return err == nil
}

func hashPassword(plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// GetStringOrDefault returns dereferenced string or default value
func GetStringOrDefault(s *string, def string) string {
	if s != nil {
		return *s
	}
	return def
}

// GetBoolOrDefault returns dereferenced bool or default value
func GetBoolOrDefault(b *bool, def bool) bool {
	if b != nil {
		return *b
	}
	return def
}

// ParseDateFromPointer parses string pointer date or returns default (zero time or error)
func ParseDateFromPointer(dateStr *string, layout string) (time.Time, error) {
	if dateStr == nil {
		return time.Time{}, nil
	}
	return time.Parse(layout, *dateStr)
}
