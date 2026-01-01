package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/container"
	"github.com/leventeberry/goapi/controllers"
)

// SetupRoutes registers all application routes on the provided Gin engine
// Uses dependency injection container for all dependencies
func SetupRoutes(router *gin.Engine, c *container.Container) {
	// Home / welcome message
	// @Summary      Welcome message
	// @Description  Returns API welcome message
	// @Tags         health
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  map[string]interface{}  "API is running"
	// @Router       / [get]
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome!",
			"status":  http.StatusOK,
		})
	})

	// Health check endpoint
	// @Summary      Health check
	// @Description  Comprehensive health check verifying database and Redis connectivity
	// @Tags         health
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  map[string]interface{}  "All systems healthy"
	// @Failure      503  {object}  map[string]interface{}  "Service unavailable"
	// @Router       /health [get]
	router.GET("/health", healthCheckHandler(c))

	// API v1 routes group
	// All API endpoints are versioned under /api/v1 for backward compatibility
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		// @Summary      Login user
		// @Description  Authenticate a user with email and password, returns JWT token
		// @Tags         authentication
		// @Accept       json
		// @Produce      json
		// @Param        credentials  body      RequestUserInput  true  "Login credentials"
		// @Success      200          {object}  map[string]interface{}  "Login successful"
		// @Failure      400          {object}  map[string]string  "Invalid request"
		// @Failure      401          {object}  map[string]string  "Invalid credentials"
		// @Failure      500          {object}  map[string]string  "Server error"
		// @Router       /api/v1/login [post]
		v1.POST("/login", controllers.LoginUser(c.AuthService))

		// @Summary      Register new user
		// @Description  Create a new user account and receive JWT token
		// @Tags         authentication
		// @Accept       json
		// @Produce      json
		// @Param        user  body      SignupUserInput  true  "User registration data"
		// @Success      200   {object}  map[string]interface{}  "Registration successful"
		// @Failure      400   {object}  map[string]string  "Invalid request"
		// @Failure      409   {object}  map[string]string  "Email already registered"
		// @Failure      500   {object}  map[string]string  "Server error"
		// @Router       /api/v1/register [post]
		v1.POST("/register", controllers.SignupUser(c.AuthService))

		// User routes setup
		SetupUserRoutes(v1, c)
	}
}

// healthCheckHandler returns a handler function for the health check endpoint
func healthCheckHandler(container *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		}

		// Check database connectivity
		sqlDB, err := container.DB.DB()
		if err != nil {
			health["database"] = gin.H{"status": "unhealthy", "error": "failed to get database connection"}
			health["status"] = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			health["database"] = gin.H{"status": "unhealthy", "error": err.Error()}
			health["status"] = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		health["database"] = gin.H{"status": "healthy"}

		// Check Redis/cache connectivity (if enabled)
		if container.Cache != nil {
			if err := container.Cache.Ping(ctx); err != nil {
				health["cache"] = gin.H{"status": "unhealthy", "error": err.Error()}
				// Cache is optional, so don't fail overall health if cache is down
				// but still report it in the response
			} else {
				health["cache"] = gin.H{"status": "healthy"}
			}
		} else {
			health["cache"] = gin.H{"status": "disabled"}
		}

		// Determine overall status code
		statusCode := http.StatusOK
		if health["status"] == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, health)
	}
}
