package utils

import (
    "golang.org/x/crypto/bcrypt"
    "errors"
)

// HashPassword hashes the password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", errors.New("failed to hash password")
    }
    return string(bytes), nil
}

// CheckPasswordHash checks the password hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}