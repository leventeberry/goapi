package services

import (
	"context"
	"errors"
	"strings"

	"github.com/leventeberry/goapi/cache"
	"github.com/leventeberry/goapi/logger"
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/models"
	"github.com/leventeberry/goapi/repositories"
)

// userService implements UserService interface
type userService struct {
	userRepo repositories.UserRepository
	cache    cache.Cache
}

// NewUserService creates a new instance of UserService
// Factory function for creating user service
func NewUserService(userRepo repositories.UserRepository, cacheClient cache.Cache) UserService {
	return &userService{
		userRepo: userRepo,
		cache:    cacheClient,
	}
}

// CreateUser creates a new user with business logic validation
func (s *userService) CreateUser(ctx context.Context, input *CreateUserInput) (*models.User, error) {
	// Validate role
	if input.Role != "" && !IsValidRole(input.Role) {
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

	// Store in cache after successful creation
	if err := s.cache.SetUserByID(ctx, user.ID, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Int("user_id", user.ID).Msg("Failed to cache user by ID")
	}
	if err := s.cache.SetUserByEmail(ctx, user.Email, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Str("email", user.Email).Msg("Failed to cache user by email")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID using cache-aside pattern
// 1. Check cache first
// 2. If cache miss, query database
// 3. Store result in cache for future requests
func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	// Try to get from cache first
	user, err := s.cache.GetUserByID(ctx, id)
	if err == nil {
		// Cache hit - return cached user
		return user, nil
	}

	// Cache miss or error - fallback to database
	if !errors.Is(err, cache.ErrCacheMiss) {
		logger.Log.Warn().Err(err).Int("user_id", id).Msg("Cache error when fetching user by ID")
	}

	user, err = s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Store in cache for future requests (best effort - don't fail on cache error)
	if err := s.cache.SetUserByID(ctx, id, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Int("user_id", id).Msg("Failed to cache user by ID")
	}
	if err := s.cache.SetUserByEmail(ctx, user.Email, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Str("email", user.Email).Msg("Failed to cache user by email")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email using cache-aside pattern
// 1. Check cache first
// 2. If cache miss, query database
// 3. Store result in cache for future requests
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Try to get from cache first
	user, err := s.cache.GetUserByEmail(ctx, email)
	if err == nil {
		// Cache hit - return cached user
		return user, nil
	}

	// Cache miss or error - fallback to database
	if !errors.Is(err, cache.ErrCacheMiss) {
		logger.Log.Warn().Err(err).Str("email", email).Msg("Cache error when fetching user by email")
	}

	user, err = s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Store in cache for future requests (best effort - don't fail on cache error)
	if err := s.cache.SetUserByEmail(ctx, email, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Str("email", email).Msg("Failed to cache user by email")
	}
	if err := s.cache.SetUserByID(ctx, user.ID, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Int("user_id", user.ID).Msg("Failed to cache user by ID")
	}

	return user, nil
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll()
}

// UpdateUser updates a user with business logic validation
func (s *userService) UpdateUser(ctx context.Context, id int, input *UpdateUserInput) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Store old email for cache invalidation if email is being changed
	oldEmail := user.Email

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
	// Compare normalized emails to handle case differences
	if input.Email != nil {
		normalizedInputEmail := strings.ToLower(strings.TrimSpace(*input.Email))
		normalizedCurrentEmail := strings.ToLower(strings.TrimSpace(user.Email))
		if normalizedInputEmail != normalizedCurrentEmail {
			exists, err := s.userRepo.ExistsByEmail(*input.Email)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrEmailExists
			}
			user.Email = *input.Email // Repository will normalize on save
		}
	}

	// Handle role update with validation
	if input.Role != nil {
		if !IsValidRole(*input.Role) {
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

	// Invalidate cache - delete old entries
	// If email changed, delete both old and new email keys
	if input.Email != nil && *input.Email != oldEmail {
		// Delete old email key
		s.cache.DeleteUserByEmail(ctx, oldEmail)
		// Delete ID key (will be repopulated on next read)
		s.cache.DeleteUserByID(ctx, id)
	} else {
		// Delete all cached entries for this user (both ID and email)
		s.cache.DeleteUser(ctx, id, user.Email)
	}

	// Store updated user in cache for future requests
	if err := s.cache.SetUserByID(ctx, user.ID, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Int("user_id", user.ID).Msg("Failed to cache updated user by ID")
	}
	if err := s.cache.SetUserByEmail(ctx, user.Email, user, cache.UserCacheTTL); err != nil {
		logger.Log.Warn().Err(err).Str("email", user.Email).Msg("Failed to cache updated user by email")
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Get user first to get email for cache invalidation
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	email := user.Email

	// Delete from database
	err = s.userRepo.Delete(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Invalidate cache - delete all cached entries for this user
	s.cache.DeleteUser(ctx, id, email)

	return nil
}

// ValidateRole checks if a role is valid
// Uses the shared IsValidRole function for consistency
func (s *userService) ValidateRole(role string) bool {
	return IsValidRole(role)
}
