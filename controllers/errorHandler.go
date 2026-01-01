package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/services"
)

// handleServiceError converts service errors to appropriate HTTP responses
func handleServiceError(c *gin.Context, err error) {
	switch err {
	case services.ErrInvalidCredentials:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	case services.ErrEmailExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
	case services.ErrInvalidRole:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Valid roles are: user, admin"})
	case services.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case services.ErrNoFieldsToUpdate:
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

