package cache

import "errors"

// Cache errors
var (
	ErrCacheMiss        = errors.New("cache miss")
	ErrCacheKeyNotFound = errors.New("cache key not found")
)
