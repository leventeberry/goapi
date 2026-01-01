package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/container"
	"github.com/leventeberry/goapi/controllers"
	"github.com/leventeberry/goapi/middleware"
)

// SetupUserRoutes registers all user-related routes on the provided Gin engine
// All user routes are protected by authentication middleware
// Admin-only routes use RequireRole middleware for role-based access control
// Uses dependency injection container for all dependencies
func SetupUserRoutes(router *gin.Engine, c *container.Container) {
	// User routes group with authentication middleware
	userGroup := router.Group("/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		// Public authenticated routes (any authenticated user can access)
		userGroup.GET("", controllers.GetUsers(c.UserService))
		userGroup.GET("/:id", controllers.GetUser(c.UserService))
		userGroup.POST("", controllers.CreateUser(c.UserService))
		userGroup.PUT("/:id", controllers.UpdateUser(c.UserService))

		// Admin-only routes (require admin role)
		userGroup.DELETE("/:id", middleware.RequireRole(c.UserRepository, "admin"), controllers.DeleteUser(c.UserService))
	}
}
