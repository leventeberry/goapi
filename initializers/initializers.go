package initializer

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
  	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func LoadEnvVariables() {
	fmt.Println("Loading Environment Variables...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Environment Variables Loaded Successfully")
}

func ConnectDatabase() {
	fmt.Println("Connecting to Database...")
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	DB = db
	fmt.Println("Connected to Database Successfully")
}
