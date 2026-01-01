package services

import (
	"errors"
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/models"
	"github.com/leventeberry/goapi/repositories"
)

// userService implements UserService interface
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
// Factory function for creating user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with business logic validation
func (s *userService) CreateUser(input *CreateUserInput) (*models.User, error) {
	// Validate role
	if input.Role != "" && !s.ValidateRole(input.Role) {
		return nil, ErrInvalidRole
	}

	// Set default role
	role := input.Role
	if role == "" {
		role = "user"
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hash, err := middleware.HashPassword(input.Password)
	if err != nil {
		return nil, ErrPasswordHashing
	}

	// Create user model
	user := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		PassHash:  hash,
		PhoneNum:  input.PhoneNum,
		Role:      role,
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// UpdateUser updates a user with business logic validation
func (s *userService) UpdateUser(id int, input *UpdateUserInput) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Validate at least one field is being updated
	if input.FirstName == nil && input.LastName == nil && input.Email == nil &&
		input.Password == nil && input.PhoneNum == nil && input.Role == nil {
		return nil, ErrNoFieldsToUpdate
	}

	// Update fields if provided
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.PhoneNum != nil {
		user.PhoneNum = *input.PhoneNum
	}

	// Handle email update with uniqueness check
	if input.Email != nil && *input.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(*input.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
		user.Email = *input.Email
	}

	// Handle role update with validation
	if input.Role != nil {
		if !s.ValidateRole(*input.Role) {
			return nil, ErrInvalidRole
		}
		user.Role = *input.Role
	}

	// Handle password update with hashing
	if input.Password != nil {
		hash, err := middleware.HashPassword(*input.Password)
		if err != nil {
			return nil, ErrPasswordHashing
		}
		user.PassHash = hash
	}

	// Save updates
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id int) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// ValidateRole checks if a role is valid
func (s *userService) ValidateRole(role string) bool {
	validRoles := []string{"user", "admin"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

