package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/leventeberry/goapi/models"
	"github.com/redis/go-redis/v9"
)

// redisCache implements Cache interface using Redis
type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache implementation
func NewRedisCache(client *redis.Client) Cache {
	if client == nil {
		return NewNoOpCache()
	}
	return &redisCache{
		client: client,
	}
}

// GetUserByID retrieves a user from cache by ID
func (r *redisCache) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	key := fmt.Sprintf("user:id:%d", id)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUserByID stores a user in cache by ID
func (r *redisCache) SetUserByID(ctx context.Context, id int, user *models.User, ttl time.Duration) error {
	key := fmt.Sprintf("user:id:%d", id)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetUserByEmail retrieves a user from cache by email
func (r *redisCache) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	key := fmt.Sprintf("user:email:%s", email)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUserByEmail stores a user in cache by email
func (r *redisCache) SetUserByEmail(ctx context.Context, email string, user *models.User, ttl time.Duration) error {
	key := fmt.Sprintf("user:email:%s", email)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

// DeleteUserByID deletes a user from cache by ID
func (r *redisCache) DeleteUserByID(ctx context.Context, id int) error {
	key := fmt.Sprintf("user:id:%d", id)
	return r.client.Del(ctx, key).Err()
}

// DeleteUserByEmail deletes a user from cache by email
func (r *redisCache) DeleteUserByEmail(ctx context.Context, email string) error {
	key := fmt.Sprintf("user:email:%s", email)
	return r.client.Del(ctx, key).Err()
}

// DeleteUser deletes both ID and email keys for a user
func (r *redisCache) DeleteUser(ctx context.Context, id int, email string) error {
	idKey := fmt.Sprintf("user:id:%d", id)
	emailKey := fmt.Sprintf("user:email:%s", email)
	return r.client.Del(ctx, idKey, emailKey).Err()
}

// IncrementRateLimit increments a rate limit counter and returns the new count
func (r *redisCache) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int, error) {
	rateLimitKey := fmt.Sprintf("ratelimit:%s", key)
	count, err := r.client.Incr(ctx, rateLimitKey).Result()
	if err != nil {
		return 0, err
	}
	
	// Set expiration on first increment
	if count == 1 {
		r.client.Expire(ctx, rateLimitKey, window)
	}
	
	return int(count), nil
}

// GetRateLimit gets the current rate limit count
func (r *redisCache) GetRateLimit(ctx context.Context, key string) (int, error) {
	rateLimitKey := fmt.Sprintf("ratelimit:%s", key)
	count, err := r.client.Get(ctx, rateLimitKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return int(count), nil
}

// ResetRateLimit resets a rate limit counter
func (r *redisCache) ResetRateLimit(ctx context.Context, key string) error {
	rateLimitKey := fmt.Sprintf("ratelimit:%s", key)
	return r.client.Del(ctx, rateLimitKey).Err()
}

// Get retrieves a value from cache by key
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrCacheMiss
		}
		return "", err
	}
	return val, nil
}

// Set stores a value in cache with TTL
func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key from cache
func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Ping checks if Redis connection is alive
func (r *redisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

