package controllers

import (
	"net/http"
	"strconv"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	ID    int    `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	PassHash string `json:"password_hash"`
	PhoneNum string `json:"phone_number"`
	Role string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdateAt string `json:"updated_at"`
}

func GetUsers(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Prepare destination slice
        var users []User

        // 2. Fetch all users; GORM populates 'users' and returns a *gorm.DB
        result := db.Find(&users)
        if result.Error != nil {
            // 3. Handle error
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
            return
        }

        // 4. Return the users slice as JSON
        c.JSON(http.StatusOK, users)
    }
}

func GetUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Parse & validate the ID
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        // 2. Attempt to load the user
        var user User
        res := db.First(&user, id)
        if res.Error != nil {
            if errors.Is(res.Error, gorm.ErrRecordNotFound) {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            }
            return
        }

        // 3. Return the user
        c.JSON(http.StatusOK, user)
    }
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Bind incoming JSON into a User struct
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }

        // 2. Let GORM insert the new record
        result := db.Create(&user)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        // 3. Return the created user (with its new ID) and 201 status
        c.JSON(http.StatusCreated, user)
    }
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Parse & validate the ID
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        // 2. Bind incoming JSON into a temporary user
        var input User
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }

        // 3. Load existing record
        var user User
        if err := db.First(&user, id).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            }
            return
        }

        // 4. Copy over fields to be updated
        user.FirstName = input.FirstName
        user.LastName  = input.LastName
        user.Email     = input.Email
        user.PassHash  = input.PassHash
        user.PhoneNum  = input.PhoneNum
        user.Role      = input.Role

        // 5. Save updates
        if err := db.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
            return
        }

        // 6. Return the updated user
        c.JSON(http.StatusOK, user)
    }
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Parse & validate the ID
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        // 2. Perform the delete
        result := db.Delete(&User{}, id)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }

        // 3. Handle not-found
        if result.RowsAffected == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        // 4. Success
        c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
    }
}