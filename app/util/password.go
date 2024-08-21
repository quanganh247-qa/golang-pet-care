package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword return the bcrypt of hash password
func HashPassword(password string) (string, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hass password")
	}
	return string(hasedPassword), nil
}

// CheckPassword checks if provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
