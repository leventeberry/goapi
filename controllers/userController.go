package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/models"
	"github.com/leventeberry/goapi/services"
)
// CreateUserInput holds the data for creating a new user
type CreateUserInput struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=128"`
	PhoneNum  string `json:"phone_number" binding:"omitempty,max=20"`
	Role      string `json:"role" binding:"omitempty,oneof=user admin"`
}

// UpdateUserInput holds the data for updating a user
type UpdateUserInput struct {
	FirstName *string `json:"first_name" binding:"omitempty,min=1,max=50"`
	LastName  *string `json:"last_name" binding:"omitempty,min=1,max=50"`
	Email     *string `json:"email" binding:"omitempty,email,max=255"`
	Password  *string `json:"password" binding:"omitempty,min=8,max=128"`
	PhoneNum  *string `json:"phone_number" binding:"omitempty,max=20"`
	Role      *string `json:"role" binding:"omitempty,oneof=user admin"`
}

// UserResponse represents a user in API responses
// Excludes sensitive fields like password hash
type UserResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhoneNum  string `json:"phone_number"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toUserResponse converts a models.User to UserResponse
func toUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		PhoneNum:  user.PhoneNum,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

// toUserResponseList converts a slice of models.User to []UserResponse
func toUserResponseList(users []models.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i := range users {
		responses[i] = *toUserResponse(&users[i])
	}
	return responses
}

// GetUsers retrieves all users with optional pagination
// @Summary      Get all users
// @Description  Get a list of all users with optional pagination (requires authentication). Query parameters: page (default: 1), page_size (default: 10, max: 100)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        page_size  query     int     false  "Items per page (default: 10, max: 100)"
// @Success      200        {object}  map[string]interface{}  "Paginated users response"
// @Failure      400        {object}  map[string]string  "Invalid pagination parameters"
// @Failure      401        {object}  map[string]string  "Unauthorized"
// @Failure      500        {object}  map[string]string  "Server error"
// @Router       /users [get]
func GetUsers(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if pagination parameters are provided
		pageParam := c.Query("page")
		pageSizeParam := c.Query("page_size")

		if pageParam != "" || pageSizeParam != "" {
			// Use pagination
			page := 1
			pageSize := 10

			if pageParam != "" {
				parsedPage, err := strconv.Atoi(pageParam)
				if err != nil || parsedPage < 1 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
					return
				}
				page = parsedPage
			}

			if pageSizeParam != "" {
				parsedPageSize, err := strconv.Atoi(pageSizeParam)
				if err != nil || parsedPageSize < 1 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size parameter"})
					return
				}
				pageSize = parsedPageSize
			}

			params := &services.PaginationParams{
				Page:     page,
				PageSize: pageSize,
			}

			users, total, err := userService.GetAllUsersPaginated(c.Request.Context(), params)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
				return
			}

			// Apply capping logic here to match service behavior
			actualPageSize := pageSize
			if actualPageSize < 1 {
				actualPageSize = 10
			}
			if actualPageSize > 100 {
				actualPageSize = 100
			}

			c.JSON(http.StatusOK, gin.H{
				"data":        users,
				"total":       total,
				"page":        page,
				"page_size":   actualPageSize,
				"total_pages": (int(total) + actualPageSize - 1) / actualPageSize, // Ceiling division
			})
			return
		}

		// No pagination parameters - return all users (backward compatibility)
		users, err := userService.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		c.JSON(http.StatusOK, toUserResponseList(users))
	}
}

// GetUser retrieves a specific user by ID
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
func GetUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := userService.GetUserByID(c.Request.Context(), int(id))
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusOK, toUserResponse(user))
	}
}

// CreateUser creates a new user
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
func CreateUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate password strength
		if err := services.ValidatePasswordStrength(input.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createInput := &services.CreateUserInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			PhoneNum:  input.PhoneNum,
			Role:      input.Role,
		}

		user, err := userService.CreateUser(c.Request.Context(), createInput)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusCreated, toUserResponse(user))
	}
}

// UpdateUser updates an existing user
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
func UpdateUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var input UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate password strength if password is being updated
		if input.Password != nil {
			if err := services.ValidatePasswordStrength(*input.Password); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		updateInput := &services.UpdateUserInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			PhoneNum:  input.PhoneNum,
			Role:      input.Role,
		}

		user, err := userService.UpdateUser(c.Request.Context(), int(id), updateInput)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusOK, toUserResponse(user))
	}
}

// DeleteUser deletes a user (admin only)
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
func DeleteUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = userService.DeleteUser(c.Request.Context(), int(id))
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
