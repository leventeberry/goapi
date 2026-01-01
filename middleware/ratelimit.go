package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
	// globalRateLimiter is a singleton rate limiter instance
	globalRateLimiter *RateLimiter
	rateLimiterOnce   sync.Once
)

// RateLimitMiddleware returns a middleware that rate limits requests per IP
// Default: 60 requests per minute with burst of 10
func RateLimitMiddleware() gin.HandlerFunc {
	// Initialize rate limiter once (singleton pattern)
	rateLimiterOnce.Do(func() {
		// Default configuration: 60 requests per minute, burst of 10
		config := RateLimiterConfig{
			RequestsPerMinute: 60,
			BurstSize:         10,
		}
		globalRateLimiter = NewRateLimiter(config)
	})

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !globalRateLimiter.allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

