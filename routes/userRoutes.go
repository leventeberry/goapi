package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/controllers"
	"github.com/leventeberry/goapi/middleware"
	"gorm.io/gorm"
)

// SetupUserRoutes registers all user-related routes on the provided Gin engine.
// All user routes are protected by authentication middleware.
func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {

	// User routes group with authentication middleware
	userGroup := router.Group("/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("", controllers.GetUsers(db))
		userGroup.GET("/:id", controllers.GetUser(db))
		userGroup.POST("", controllers.CreateUser(db))
		userGroup.PUT("/:id", controllers.UpdateUser(db))
		userGroup.DELETE("/:id", controllers.DeleteUser(db))
	}
}
