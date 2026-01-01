package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/logger"
)

// RequestLogger returns a middleware that logs HTTP requests with details.
// Logs: method, path, status code, response time, client IP, and user agent.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Get status code
		statusCode := c.Writer.Status()

		// Get method
		method := c.Request.Method

		// Build full path with query string if present
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		// Create structured log event with common fields
		baseEvent := logger.Log.
			With().
			Str("method", method).
			Str("path", fullPath).
			Str("proto", c.Request.Proto).
			Int("status_code", statusCode).
			Dur("latency", latency).
			Str("client_ip", clientIP).
			Str("user_agent", userAgent).
			Logger()

		// Log based on status code with appropriate level
		if statusCode >= 500 {
			baseEvent.Error().Msg("HTTP Request")
		} else if statusCode >= 400 {
			baseEvent.Warn().Msg("HTTP Request")
		} else {
			baseEvent.Info().Msg("HTTP Request")
		}
	}
}
