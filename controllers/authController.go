package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/middleware"
)

type ExistUser struct {
	ID        int    `json:"id"`
}

type RequestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	PhoneNum  string `json:"phone_number"`
	Role      string `json:"role"`
}

func LoginUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Bind the request body to the User struct
		var user RequestUser
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	// Query the database
		row := db.QueryRow("SELECT * FROM users WHERE email = ?", user.Email)

		// Scan the row into a User struct
		var dbUser User
		err := row.Scan(
			&dbUser.ID,
			&dbUser.FirstName,
			&dbUser.LastName,
			&dbUser.Email,
			&dbUser.PassHash,
			&dbUser.PhoneNum,
			&dbUser.Role,
			&dbUser.CreatedAt,
			&dbUser.UpdateAt,
		)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user from the database"})
			return
		}

		// Compare the password with the stored hash
		if !middleware.ComparePasswords(dbUser.PassHash, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		// Create a new JWT token
		token, err := middleware.CreateToken(dbUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		// Return the token
		c.JSON(http.StatusOK, gin.H{"token": token})

	}
}

func SignupUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request body to the User struct
		var user RegisterUser
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user already exists
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&exists)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}


		// Hash the password
		passHash, err := middleware.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Insert the user into the database
		result, err := db.Exec("INSERT INTO users (email, password_hash, first_name, last_name, phone_number, role) VALUES (?, ?, ?, ?, ?, ?)", user.Email, passHash, user.FirstName, user.LastName, user.PhoneNum, user.Role)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user into the database"})
			return
		}

		// Get the ID of the inserted user
		userID, err := result.LastInsertId()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
			return
		}

		// Create a new JWT token
		token, err := middleware.CreateToken(int(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		// Return the token
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id": userID,
				"email": user.Email,
			},
		})
	}
}
