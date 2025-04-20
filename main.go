package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/leventeberry/goapi/controllers"
    "github.com/leventeberry/goapi/initializers"
)

func main() {
    // Initialize environment variables and database connection, then run migrations
    initializers.Init()

    // Create a Gin router with default middleware (logger and recovery)
    router := gin.Default()

    // Home route / health check
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to the API",
            "status":  http.StatusOK,
        })
    })

    // Authentication routes
    router.POST("/login", controllers.LoginUser(initializers.DB))
    router.POST("/register", controllers.SignupUser(initializers.DB))

    // User CRUD routes grouped under /users
    userGroup := router.Group("/users")
    {
        userGroup.GET("", controllers.GetUsers(initializers.DB))
        userGroup.GET("/:id", controllers.GetUser(initializers.DB))
        userGroup.POST("", controllers.CreateUser(initializers.DB))
        userGroup.PUT("/:id", controllers.UpdateUser(initializers.DB))
        userGroup.DELETE("/:id", controllers.DeleteUser(initializers.DB))
    }

    // Start server on specified PORT or default to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}