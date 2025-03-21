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

    // User Routes
    router.GET("/users", controllers.GetUsers(db))
    router.GET("/users/:id", controllers.GetUser(db))
    router.POST("/users", controllers.CreateUser(db))
    router.PUT("/users/:id", controllers.UpdateUser(db))
    router.DELETE("/users/:id", controllers.DeleteUser(db))

    // Start the server (blocking call)
    router.Run(":8080")
}

