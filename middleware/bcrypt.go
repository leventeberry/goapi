package middleware

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

// PasswordCost defines the bcrypt hashing cost.
// Increase this if you need stronger hashes at the expense of CPU time.
const PasswordCost = bcrypt.DefaultCost

// HashPassword generates a bcrypt hash of the given plaintext password.
func HashPassword(password string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
    if err != nil {
        return "", fmt.Errorf("password hashing failed: %w", err)
    }
    return string(hashed), nil
}

// ComparePasswords checks whether the provided plaintext password matches
// the stored bcrypt hash. Returns true if they match.
func ComparePasswords(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

