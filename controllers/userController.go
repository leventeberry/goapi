package controllers

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	// Query the database
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Failed to query the database"})
		return
	}

	// Iterate over the rows
	var users []map[string]interface{}
	for rows.Next() {
		// Create a new map to store the column name and value
		user := make(map[string]interface{})
		// Scan the row into the map
		err = rows.ScanMap(user)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": "Failed to scan the row"})
			return
		}
		// Append the map to the users slice
		users = append(users, user)
	}
	}
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var user map[string]interface{}
		err := c.BindJSON(&user)
		if err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": "Failed to parse the request body"})
			return
		}

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user["name"], user["email"])
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": "Failed to insert the user into the database"})
			return
		}

		// Return a success response
		c.JSON(200, gin.H{"message": "User created successfully"})
	}
}