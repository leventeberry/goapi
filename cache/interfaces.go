package cache

import (
	"context"
	"time"

	"github.com/leventeberry/goapi/models"
)

// Cache defines the interface for cache operations
// Supports both user caching and rate limiting operations
type Cache interface {
	// User cache operations
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	SetUserByID(ctx context.Context, id int, user *models.User, ttl time.Duration) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SetUserByEmail(ctx context.Context, email string, user *models.User, ttl time.Duration) error
	DeleteUserByID(ctx context.Context, id int) error
	DeleteUserByEmail(ctx context.Context, email string) error
	DeleteUser(ctx context.Context, id int, email string) error // Deletes both ID and email keys

	// Rate limiting operations
	IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int, error)
	GetRateLimit(ctx context.Context, key string) (int, error)
	ResetRateLimit(ctx context.Context, key string) error

	// General cache operations
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Health check
	Ping(ctx context.Context) error
}
