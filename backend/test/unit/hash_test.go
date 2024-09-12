package unit_test

import (
	"testing"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
)

// TestHashPassword tests the HashPassword function
func TestHashPassword(t *testing.T) {
	password := "testpassword"

	// Test successful password hashing
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(hashedPassword) == 0 {
		t.Fatalf("Expected hashed password to be non-empty")
	}

	// Test hashing with empty password
	_, err = security.HashPassword("")
	if err == nil {
		t.Fatalf("Expected an error when hashing an empty password")
	}
}

// TestCheckPasswordHash tests the CheckPasswordHash function
func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"

	// Hash the password first
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test valid password check
	if !security.CheckPasswordHash(password, hashedPassword) {
		t.Fatalf("Expected CheckPasswordHash to return true for a valid password and hash")
	}

	// Test invalid password check
	if security.CheckPasswordHash("wrongpassword", hashedPassword) {
		t.Fatalf("Expected CheckPasswordHash to return false for an invalid password and hash")
	}

	// Test invalid hash format
	if security.CheckPasswordHash(password, "invalidhash") {
		t.Fatalf("Expected CheckPasswordHash to return false for an invalid hash format")
	}
}
