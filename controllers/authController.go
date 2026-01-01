package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/models"
	"gorm.io/gorm"
)

// Valid roles that can be assigned to users
var ValidRoles = []string{"user", "admin"}

// isValidRole checks if the provided role is one of the valid roles
func isValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// RequestUser holds login credentials.
type RequestUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignupUserInput holds registration data.
type SignupUserInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	PhoneNum  string `json:"phone_number" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

// Function that returns data upon successful query
func ReturnSuccessData(c *gin.Context, user *models.User, token *middleware.Authentication) {
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})

}

// LoginUser authenticates a user and returns a JWT token.
// @Summary      Login user
// @Description  Authenticate a user with email and password, returns JWT token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      RequestUserInput  true  "Login credentials"
// @Success      200          {object}  map[string]interface{}  "Login successful"
// @Failure      400          {object}  map[string]string  "Invalid request"
// @Failure      401          {object}  map[string]string  "Invalid credentials"
// @Failure      500          {object}  map[string]string  "Server error"
// @Router       /login [post]
func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RequestUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Load user by email
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// Check password
		if !middleware.ComparePasswords(user.PassHash, input.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate JWT
		token, err := middleware.CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		ReturnSuccessData(c, &user, token)
	}
}

// SignupUser registers a new user and returns a JWT token.
// @Summary      Register new user
// @Description  Create a new user account and receive JWT token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        user  body      SignupUserInput  true  "User registration data"
// @Success      200   {object}  map[string]interface{}  "Registration successful"
// @Failure      400   {object}  map[string]string  "Invalid request"
// @Failure      409   {object}  map[string]string  "Email already registered"
// @Failure      500   {object}  map[string]string  "Server error"
// @Router       /register [post]
func SignupUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignupUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if email already exists
		var existing models.User
		if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Validate role
		if !isValidRole(input.Role) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Valid roles are: user, admin"})
			return
		}

		// Hash password
		hash, err := middleware.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Create user record
		user := models.User{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			PassHash:  hash,
			PhoneNum:  input.PhoneNum,
			Role:      input.Role,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Generate JWT
		token, err := middleware.CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		ReturnSuccessData(c, &user, token)
	}
}
