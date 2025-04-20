package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/leventeberry/goapi/controllers"
    "gorm.io/gorm"
)

// SetupRoutes registers all application routes on the provided Gin engine.
func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {

    // User routes group
    userGroup := router.Group("/users")
    {
        userGroup.GET("", controllers.GetUsers(db))
        userGroup.GET("/:id", controllers.GetUser(db))
        userGroup.POST("", controllers.CreateUser(db))
        userGroup.PUT("/:id", controllers.UpdateUser(db))
        userGroup.DELETE("/:id", controllers.DeleteUser(db))
    }
}