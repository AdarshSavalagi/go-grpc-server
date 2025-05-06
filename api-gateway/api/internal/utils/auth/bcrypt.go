package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptCost defines the cost factor for hashing. Higher = more secure but slower.
const bcryptCost = bcrypt.DefaultCost // Can be replaced with a custom value (e.g., 12)

// HashPassword securely hashes a plaintext password using bcrypt.
//
// Parameters:
//   - password: The plain text password to be hashed.
//
// Returns:
//   - The hashed password as a string.
//   - An error if the hashing process fails.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a hashed password with a plain text password.
//
// Parameters:
//   - hashedPassword: The previously hashed password.
//   - password: The plain text password to verify.
//
// Returns:
//   - true if the password matches the hash, false otherwise.
func VerifyPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
