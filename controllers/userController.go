package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/middleware"
	"github.com/leventeberry/goapi/models"
	"gorm.io/gorm"
)

// GetUsers retrieves all users.
// @Summary      Get all users
// @Description  Get a list of all users (requires authentication)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.User  "List of users"
// @Failure      401  {object}  map[string]string  "Unauthorized"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Prepare destination slice
		var users []models.User

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

// GetUser retrieves a specific user by ID.
// @Summary      Get user by ID
// @Description  Get a specific user by their ID (requires authentication)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User  "User object"
// @Failure      400  {object}  map[string]string  "Invalid user ID"
// @Failure      401  {object}  map[string]string  "Unauthorized"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /users/{id} [get]
func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Parse & validate the ID
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// 2. Attempt to load the user
		var user models.User
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

// CreateUserInput holds the data for creating a new user.
type CreateUserInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	PhoneNum  string `json:"phone_number"`
	Role      string `json:"role"`
}

// CreateUser creates a new user.
// @Summary      Create new user
// @Description  Create a new user account (requires authentication)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user  body      CreateUserInput  true  "User data"
// @Success      201   {object}  models.User  "Created user"
// @Failure      400   {object}  map[string]string  "Invalid request"
// @Failure      401   {object}  map[string]string  "Unauthorized"
// @Failure      409   {object}  map[string]string  "Email already registered"
// @Failure      500   {object}  map[string]string  "Server error"
// @Router       /users [post]
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Bind incoming JSON into input struct
		var input CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// 2. Check if email already exists
		var existing models.User
		if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 3. Hash password
		hash, err := middleware.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// 4. Set default role if not provided and validate role
		role := input.Role
		if role == "" {
			role = "user"
		} else if !isValidRole(role) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Valid roles are: user, admin"})
			return
		}

		// 5. Create user record
		user := models.User{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			PassHash:  hash,
			PhoneNum:  input.PhoneNum,
			Role:      role,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// 6. Return the created user (with its new ID) and 201 status
		c.JSON(http.StatusCreated, user)
	}
}

// UpdateUserInput holds the data for updating a user.
type UpdateUserInput struct {
	FirstName string `json:"first_name" binding:"omitempty,min=1,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=1,max=50"`
	Email     string `json:"email" binding:"omitempty,email"`
	Password  string `json:"password" binding:"omitempty,min=8"`
	PhoneNum  string `json:"phone_number" binding:"omitempty,max=20"`
	Role      string `json:"role" binding:"omitempty"`
}

// UpdateUser updates an existing user.
// @Summary      Update user
// @Description  Update user information (partial updates supported, requires authentication)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int              true  "User ID"
// @Param        user  body      UpdateUserInput  true  "User update data"
// @Success      200   {object}  models.User  "Updated user"
// @Failure      400   {object}  map[string]string  "Invalid request"
// @Failure      401   {object}  map[string]string  "Unauthorized"
// @Failure      404   {object}  map[string]string  "User not found"
// @Failure      409   {object}  map[string]string  "Email already registered"
// @Failure      500   {object}  map[string]string  "Server error"
// @Router       /users/{id} [put]
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Parse & validate the ID
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// 2. Bind incoming JSON into input struct
		var input UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2a. Validate that at least one field is being updated
		if input.FirstName == "" && input.LastName == "" && input.Email == "" &&
			input.Password == "" && input.PhoneNum == "" && input.Role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
			return
		}

		// 3. Load existing record
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// 4. Check email uniqueness if email is being updated
		if input.Email != "" && input.Email != user.Email {
			var existing models.User
			if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
				return
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
			user.Email = input.Email
		}

		// 5. Update fields if provided
		if input.FirstName != "" {
			user.FirstName = input.FirstName
		}
		if input.LastName != "" {
			user.LastName = input.LastName
		}
		if input.PhoneNum != "" {
			user.PhoneNum = input.PhoneNum
		}
		if input.Role != "" {
			if !isValidRole(input.Role) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Valid roles are: user, admin"})
				return
			}
			user.Role = input.Role
		}

		// 6. Hash and update password if provided
		if input.Password != "" {
			hash, err := middleware.HashPassword(input.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			user.PassHash = hash
		}

		// 7. Save updates
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		// 8. Return the updated user
		c.JSON(http.StatusOK, user)
	}
}

// DeleteUser deletes a user (admin only).
// @Summary      Delete user
// @Description  Delete a user by ID (requires admin role)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string  "User deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid user ID"
// @Failure      401  {object}  map[string]string  "Unauthorized"
// @Failure      403  {object}  map[string]string  "Insufficient permissions"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /users/{id} [delete]
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Parse & validate the ID
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// 2. Perform the delete
		result := db.Delete(&models.User{}, id)
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
