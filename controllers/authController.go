package controllers

import (
    "net/http"
    "errors"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/leventeberry/goapi/middleware"
)

// RequestUser holds login credentials.
type RequestUser struct {
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

// LoginUser authenticates a user and returns a JWT token.
func LoginUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input RequestUser
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Load user by email
        var user User
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

        c.JSON(http.StatusOK, gin.H{"token": token})
    }
}

// SignupUser registers a new user and returns a JWT token.
func SignupUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input SignupUserInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Check if email already exists
        var existing User
        if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
            c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
            return
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }

        // Hash password
        hash, err := middleware.HashPassword(input.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
            return
        }

        // Create user record
        user := User{
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

        c.JSON(http.StatusCreated, gin.H{
            "token": token,
            "user": gin.H{
                "id":    user.ID,
                "email": user.Email,
            },
        })
    }
}
