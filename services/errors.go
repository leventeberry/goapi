package services

import "errors"

// Service errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRole        = errors.New("invalid role")
	ErrPasswordHashing    = errors.New("failed to hash password")
	ErrNoFieldsToUpdate   = errors.New("at least one field must be provided for update")
	ErrTokenGeneration    = errors.New("failed to generate token")
)

