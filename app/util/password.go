package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// PasswordConfig defines the rules for password generation
type PasswordConfig struct {
	Length        int  // Total length of password
	IncludeUpper  bool // Include uppercase letters
	IncludeLower  bool // Include lowercase letters
	IncludeNumber bool // Include numbers
	IncludeSymbol bool // Include special characters
}

// DefaultPasswordConfig returns a default secure password configuration
func DefaultPasswordConfig() PasswordConfig {
	return PasswordConfig{
		Length:        16,
		IncludeUpper:  true,
		IncludeLower:  true,
		IncludeNumber: true,
		IncludeSymbol: true,
	}
}

// GeneratePassword generates a random password based on the provided configuration
func GeneratePassword(config PasswordConfig) (string, error) {
	var (
		lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numberSet      = "0123456789"
		specialCharSet = "!@#$%&*"
		allCharSet     = ""
		password       = ""
	)

	// Build the character set based on configuration
	if config.IncludeLower {
		allCharSet += lowerCharSet
		password += getRandomChar(lowerCharSet) // Ensure at least one lowercase
	}
	if config.IncludeUpper {
		allCharSet += upperCharSet
		password += getRandomChar(upperCharSet) // Ensure at least one uppercase
	}
	if config.IncludeNumber {
		allCharSet += numberSet
		password += getRandomChar(numberSet) // Ensure at least one number
	}
	if config.IncludeSymbol {
		allCharSet += specialCharSet
		password += getRandomChar(specialCharSet) // Ensure at least one special char
	}

	// Fill the remaining length with random characters from all allowed sets
	remainingLength := config.Length - len(password)
	for i := 0; i < remainingLength; i++ {
		password += getRandomChar(allCharSet)
	}

	// Shuffle the password to make it more random
	password = shuffleString(password)

	return password, nil
}

// getRandomChar returns a random character from the provided character set
func getRandomChar(charSet string) string {
	max := big.NewInt(int64(len(charSet)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return string(charSet[0])
	}
	return string(charSet[n.Int64()])
}

// shuffleString randomly shuffles the characters in a string
func shuffleString(str string) string {
	chars := strings.Split(str, "")
	length := len(chars)

	for i := length - 1; i > 0; i-- {
		if j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1))); err == nil {
			chars[i], chars[j.Int64()] = chars[j.Int64()], chars[i]
		}
	}

	return strings.Join(chars, "")
}

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
