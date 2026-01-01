package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/models"
	"github.com/leventeberry/goapi/services"
)

// RequestUserInput holds login credentials
type RequestUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignupUserInput holds registration data
type SignupUserInput struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=128"`
	PhoneNum  string `json:"phone_number" binding:"omitempty,max=20"`
	Role      string `json:"role" binding:"omitempty,oneof=user admin"`
}

// ReturnSuccessData returns standardized success response with token and user
func ReturnSuccessData(c *gin.Context, user *models.User, token *middleware.Authentication) {
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// LoginUser authenticates a user and returns a JWT token
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
func LoginUser(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RequestUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, token, err := authService.Login(input.Email, input.Password)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		ReturnSuccessData(c, user, token)
	}
}

// SignupUser registers a new user and returns a JWT token
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
func SignupUser(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignupUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		registerInput := &services.RegisterInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			PhoneNum:  input.PhoneNum,
			Role:      input.Role,
		}

		user, token, err := authService.Register(registerInput)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		ReturnSuccessData(c, user, token)
	}

}
