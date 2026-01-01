package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/cache"
)

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// rateLimiterEntry tracks requests for a single IP
type rateLimiterEntry struct {
	tokens     int
	lastUpdate time.Time
	mu         sync.Mutex
}

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	config      RateLimiterConfig
	entries     map[string]*rateLimiterEntry
	mu          sync.RWMutex
	cleanupTick *time.Ticker
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		config:  config,
		entries: make(map[string]*rateLimiterEntry),
	}

	// Start cleanup goroutine to remove old entries
	rl.cleanupTick = time.NewTicker(5 * time.Minute)
	go rl.cleanup()

	return rl
}

// cleanup removes entries that haven't been used in the last 10 minutes
func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTick.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, entry := range rl.entries {
			entry.mu.Lock()
			if now.Sub(entry.lastUpdate) > 10*time.Minute {
				delete(rl.entries, ip)
			}
			entry.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	entry, exists := rl.entries[ip]
	if !exists {
		entry = &rateLimiterEntry{
			tokens:     rl.config.BurstSize,
			lastUpdate: time.Now(),
		}
		rl.entries[ip] = entry
	}
	rl.mu.Unlock()

	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(entry.lastUpdate)

	// Refill tokens based on time elapsed
	tokensToAdd := int(elapsed.Minutes() * float64(rl.config.RequestsPerMinute))
	if tokensToAdd > 0 {
		entry.tokens = min(entry.tokens+tokensToAdd, rl.config.BurstSize)
		entry.lastUpdate = now
	}

	// Check if we have tokens available
	if entry.tokens > 0 {
		entry.tokens--
		return true
	}

	return false
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	// globalRateLimiter is a singleton in-memory rate limiter instance
	globalRateLimiter *RateLimiter
	// globalRedisRateLimiter is a Redis-based rate limiter instance
	globalRedisRateLimiter *RedisRateLimiter
	rateLimiterOnce         sync.Once
	useRedis                bool
)

// RateLimitMiddleware returns a middleware that rate limits requests per IP
// Default: 60 requests per minute with burst of 10
// Uses in-memory rate limiting by default
func RateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddlewareWithCache(nil)
}

// RateLimitMiddlewareWithCache returns a middleware that rate limits requests per IP
// If cacheClient is provided and not a no-op cache, uses Redis-based rate limiting
// Otherwise falls back to in-memory rate limiting
// Default: 60 requests per minute with burst of 10
func RateLimitMiddlewareWithCache(cacheClient cache.Cache) gin.HandlerFunc {
	// Initialize rate limiter once (singleton pattern)
	rateLimiterOnce.Do(func() {
		// Default configuration: 60 requests per minute, burst of 10
		config := RateLimiterConfig{
			RequestsPerMinute: 60,
			BurstSize:         10,
		}

		// Check if we should use Redis
		// Use Redis if cache is provided and it's actually working (not no-op)
		if cacheClient != nil {
			// Test if cache is actually working by trying to set/get a test value
			// No-op cache will return cache miss, Redis will work
			ctx := context.Background()
			testKey := "ratelimit:init:test"
			testValue := "test"
			
			// Try to set and get - if this works, we have a real cache
			err := cacheClient.Set(ctx, testKey, testValue, time.Second)
			if err == nil {
				val, err := cacheClient.Get(ctx, testKey)
				if err == nil && val == testValue {
					// Cache is working (Redis is available), use Redis rate limiter
					globalRedisRateLimiter = NewRedisRateLimiter(cacheClient, config)
					useRedis = true
					// Clean up test key
					cacheClient.Delete(ctx, testKey)
				} else {
					// Cache not working (no-op cache), use in-memory
					globalRateLimiter = NewRateLimiter(config)
					useRedis = false
				}
			} else {
				// Cache not working, use in-memory
				globalRateLimiter = NewRateLimiter(config)
				useRedis = false
			}
		} else {
			// No cache provided, use in-memory
			globalRateLimiter = NewRateLimiter(config)
			useRedis = false
		}
	})

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		ctx := context.Background()

		var allowed bool
		if useRedis && globalRedisRateLimiter != nil {
			// Use Redis-based rate limiting
			allowed = globalRedisRateLimiter.allow(ctx, clientIP)
		} else {
			// Use in-memory rate limiting
			allowed = globalRateLimiter.allow(clientIP)
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

