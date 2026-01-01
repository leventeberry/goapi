package cache

import "time"

// Cache key patterns
const (
	// UserIDKeyPrefix is the prefix for user cache keys by ID
	// Full key format: "user:id:{id}"
	UserIDKeyPrefix = "user:id:"
	
	// UserEmailKeyPrefix is the prefix for user cache keys by email
	// Full key format: "user:email:{email}"
	UserEmailKeyPrefix = "user:email:"
	
	// RateLimitKeyPrefix is the prefix for rate limiting keys
	// Full key format: "ratelimit:{key}"
	RateLimitKeyPrefix = "ratelimit:"
)

// Cache TTL (Time To Live) values
const (
	// UserCacheTTL is the default TTL for cached user objects
	// Set to 15 minutes - balances freshness with cache efficiency
	// User data changes infrequently, so 15 minutes reduces database load
	// while ensuring reasonable data freshness
	UserCacheTTL = 15 * time.Minute
	
	// RateLimitWindow is the default window for rate limiting
	// Set to 1 minute to match typical rate limiting requirements
	// This matches the default rate limiter configuration
	RateLimitWindow = 1 * time.Minute
)

