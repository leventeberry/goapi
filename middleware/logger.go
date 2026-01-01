package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
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

		// Build log message
		if raw != "" {
			path = path + "?" + raw
		}

		logMessage := fmt.Sprintf(
			"[%s] %s %s %d %v %s %s",
			method,
			path,
			c.Request.Proto,
			statusCode,
			latency,
			clientIP,
			userAgent,
		)

		// Log based on status code
		if statusCode >= 500 {
			log.Printf("ERROR: %s", logMessage)
		} else if statusCode >= 400 {
			log.Printf("WARN: %s", logMessage)
		} else {
			log.Printf("INFO: %s", logMessage)
		}
	}
}

