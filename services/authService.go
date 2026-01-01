package services

import (
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/models"
	"github.com/leventeberry/goapi/repositories"
)

// authService implements AuthService interface
type authService struct {
	userRepo repositories.UserRepository
}

// NewAuthService creates a new instance of AuthService
// Factory function for creating auth service
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(email, password string) (*models.User, *middleware.Authentication, error) {
	// Validate credentials
	user, err := s.ValidateCredentials(email, password)
	if err != nil {
		return nil, nil, err
	}

	// Generate JWT token with user role
	token, err := middleware.CreateToken(user.ID, user.Role)
	if err != nil {
		return nil, nil, ErrTokenGeneration
	}

	return user, token, nil
}

// Register creates a new user account and returns a JWT token
func (s *authService) Register(input *RegisterInput) (*models.User, *middleware.Authentication, error) {
	// Create user directly here to avoid circular dependency
	// In a more advanced setup, we'd use a service orchestrator or composition

	// Validate role
	if input.Role != "" && !IsValidRole(input.Role) {
		return nil, nil, ErrInvalidRole
	}

	// Set default role
	role := input.Role
	if role == "" {
		role = "user"
	}

	// Check if email exists
	exists, err := s.userRepo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrEmailExists
	}

	// Hash password
	hash, err := middleware.HashPassword(input.Password)
	if err != nil {
		return nil, nil, ErrPasswordHashing
	}

	// Create user
	user := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		PassHash:  hash,
		PhoneNum:  input.PhoneNum,
		Role:      role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Generate JWT token with user role
	token, err := middleware.CreateToken(user.ID, user.Role)
	if err != nil {
		return nil, nil, ErrTokenGeneration
	}

	return user, token, nil
}

// ValidateCredentials validates user email and password
func (s *authService) ValidateCredentials(email, password string) (*models.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !middleware.ComparePasswords(user.PassHash, password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

