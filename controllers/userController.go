package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query the database
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query the database"})
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User
			err := rows.Scan(
				&user.ID,
				&user.FirstName,
				&user.LastName,
				&user.Email,
				&user.PassHash,
				&user.PhoneNum,
				&user.Role,
				&user.CreatedAt,
				&user.UpdateAt,
			)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan the row"})
				return
			}
			users = append(users, user)
		}

		// Return JSON response
		c.JSON(http.StatusOK, users)
	}
}

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the URL
		userID := c.Param("id")
		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Query the database
		row := db.QueryRow("SELECT * FROM users WHERE user_id = ?", id)

		var user User
		err = row.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PassHash,
			&user.PhoneNum,
			&user.Role,
			&user.CreatedAt,
			&user.UpdateAt,
		)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan the row"})
			return
		}

		// Return JSON response
		c.JSON(http.StatusOK, user)
	}
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var user User

		err := c.BindJSON(&user)
		if err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": "Failed to parse the request body"})
			return
		}

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO users (first_name, last_name, email, password_hash, phone_number, role) VALUES (?, ?, ?, ?, ?, ?)", &user.FirstName, &user.LastName, &user.Email, &user.PassHash, &user.PhoneNum, &user.Role)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": "Failed to insert the user into the database"})
			return
		}

		// Return a success response
		c.JSON(200, gin.H{"message": "User created successfully"})
	}
}