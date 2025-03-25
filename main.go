package main

import(
	// "fmt"
    "os"
	"github.com/gin-gonic/gin"
    "github.com/leventeberry/goapi/initializers"
	// "github.com/leventeberry/gomysql"
	// "github.com/leventeberry/goapi/controllers"
)

func init() {
    initializer.LoadEnvVariables()
    initializer.ConnectDatabase()
}

func main() {
    // Create a connection pool to the database
    // db, err := gomysql.ConnectDB("tcp", "localhost:3306", "testapi")
    // if err != nil {
    //     fmt.Println("Database connection error:", err)
    //     return
    // }

    // Create a new router
    router := gin.Default()

    // Home Route
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to the API",
            "status": 200,
        })
    })

    // router.POST("/login", controllers.LoginUser(db))
    // router.POST("/register", controllers.SignupUser(db))

    // // User Routes
    // router.GET("/users", controllers.GetUsers(db))
    // router.GET("/users/:id", controllers.GetUser(db))
    // router.POST("/users", controllers.CreateUser(db))
    // router.PUT("/users/:id", controllers.UpdateUser(db))
    // router.DELETE("/users/:id", controllers.DeleteUser(db))

    // Start the server (blocking call)
    router.Run(os.Getenv("PORT"))
}

