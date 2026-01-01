package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/leventeberry/goapi/docs"
    "github.com/leventeberry/goapi/initializers"
    "github.com/leventeberry/goapi/middleware"
    "github.com/leventeberry/goapi/routes"
)

// @title           GoAPI - RESTful API Template
// @version         1.0
// @description     A RESTful API built with Go (Golang) using the Gin web framework. This API provides user management functionality with JWT-based authentication, role-based access control (RBAC), and comprehensive middleware for security and logging.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

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

    // Create a Gin router
    router := gin.New()

    // Add middleware: rate limiter, request logger, and recovery
    router.Use(middleware.RateLimitMiddleware())
    router.Use(middleware.RequestLogger())
    router.Use(gin.Recovery())

    // Register all routes
    routes.SetupRoutes(router, initializers.DB)

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