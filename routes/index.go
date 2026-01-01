package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/container"
	"github.com/leventeberry/goapi/controllers"
)

// SetupRoutes registers all application routes on the provided Gin engine
// Uses dependency injection container for all dependencies
func SetupRoutes(router *gin.Engine, c *container.Container) {
	// Home / health check
	// @Summary      Health check
	// @Description  Returns API status and welcome message
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

	// Authentication routes
	router.POST("/login", controllers.LoginUser(c.AuthService))
	router.POST("/register", controllers.SignupUser(c.AuthService))

	// User routes setup
	SetupUserRoutes(router, c)
}
