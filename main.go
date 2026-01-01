package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/cache"
	"github.com/leventeberry/goapi/container"
	"github.com/leventeberry/goapi/docs"
	"github.com/leventeberry/goapi/initializers"
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
	var cacheClient cache.Cache
	if initializers.RedisClient != nil {
		cacheClient = cache.NewRedisCache(initializers.RedisClient)
	} else {
		cacheClient = cache.NewNoOpCache()
	}

	// Create dependency injection container using Factory Pattern
	// This initializes all repositories, services, and their dependencies
	appContainer := container.NewContainer(initializers.DB, cacheClient)

	// Create a Gin router
	router := gin.New()

	// Add middleware: rate limiter, request logger, and recovery
	router.Use(middleware.RateLimitMiddleware())
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
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
