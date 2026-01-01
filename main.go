package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/container"
	"github.com/leventeberry/goapi/docs"
	"github.com/leventeberry/goapi/initializers"
	"github.com/leventeberry/goapi/logger"
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           GoAPI - RESTful API Template
// @version         1.0
// @description     A RESTful API built with Go (Golang) using the Gin web framework. This API provides user management functionality with JWT-based authentication, role-based access control (RBAC), and comprehensive middleware for security and logging.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize Swagger docs
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Initialize environment variables, database connection, and run migrations
	initializers.Init()

	// Initialize cache client (Redis or no-op)
	// Uses helper function from initializers to centralize cache creation logic
	cacheClient := initializers.GetCacheClient()

	// Create dependency injection container using Factory Pattern
	// This initializes all repositories, services, and their dependencies
	appContainer := container.NewContainer(initializers.DB, cacheClient)

	// Create a Gin router
	router := gin.New()

	// Add middleware: rate limiter, request logger, and recovery
	// Rate limiter uses Redis if available, otherwise falls back to in-memory
	router.Use(middleware.RateLimitMiddlewareWithCache(cacheClient))
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())

	// Register all routes with dependency injection
	routes.SetupRoutes(router, appContainer)

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server on specified PORT or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup graceful shutdown
	// Listen for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := router.Run(":" + port); err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	logger.Log.Info().Str("port", port).Msg("Server is running")

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	logger.Log.Info().Msg("Shutting down server...")

	// Cleanup: close Redis connection if it exists
	initializers.CloseRedis()

	logger.Log.Info().Msg("Server exited")
}
