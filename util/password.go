package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword will hash the given password and return hashedpassword/empty string and error/nil
func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %w", err)
	}
	return string(passwordHash), nil
}

// CheckPassword will compare password to hashed password and returns error/nil
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
