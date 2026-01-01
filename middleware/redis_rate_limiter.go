package middleware

import (
	"context"
	"time"

	"github.com/leventeberry/goapi/cache"
)

// RedisRateLimiter implements rate limiting using Redis
// Uses a sliding window counter approach suitable for distributed systems
type RedisRateLimiter struct {
	cache            cache.Cache
	requestsPerMinute int
	burstSize         int
	window            time.Duration
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(cacheClient cache.Cache, config RateLimiterConfig) *RedisRateLimiter {
	return &RedisRateLimiter{
		cache:             cacheClient,
		requestsPerMinute: config.RequestsPerMinute,
		burstSize:         config.BurstSize,
		window:            cache.RateLimitWindow,
	}
}

// allow checks if a request from the given key should be allowed
// Uses Redis INCR with expiration for distributed rate limiting
func (r *RedisRateLimiter) allow(ctx context.Context, key string) bool {
	// Increment the counter for this key
	count, err := r.cache.IncrementRateLimit(ctx, key, r.window)
	if err != nil {
		// If Redis fails, allow the request (fail open)
		// In production, you might want to log this
		return true
	}

	// Check if count exceeds the limit
	// We use requestsPerMinute as the limit
	if count > r.requestsPerMinute {
		return false
	}

	return true
}

