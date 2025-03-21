package main

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leventeberry/go-mysql"
	"github.com/leventeberry/goapi/controllers"
)

func main() {
    // Create a connection pool to the database
    db, err := gomysql.ConnectDB("tcp", "localhost:3306", "testapi")
    if err != nil {
        fmt.Println("Database connection error:", err)
        return
    }

    // Create a new router
    router := gin.Default()

    // Pass handler functions correctly
    router.GET("/users", controllers.GetUsers(db))
    router.POST("/users", controllers.CreateUser(db))

    // Start the server (blocking call)
    router.Run(":8080")
}

