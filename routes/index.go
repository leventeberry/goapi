package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/leventeberry/goapi/controllers"
    "gorm.io/gorm"
)

// SetupRoutes registers all application routes on the provided Gin engine.
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    // Home / health check
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome!",
            "status":  http.StatusOK,
        })
    })

    // Authentication routes
    router.POST("/login", controllers.LoginUser(db))
    router.POST("/register", controllers.SignupUser(db))

    // User routes setup
	SetupUserRoutes(router, db)

	
}