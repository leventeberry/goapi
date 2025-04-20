package initializers

import (
    "fmt"
    "log"
    "os"
    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/leventeberry/goapi/controllers"
)

// DB is the global database connection
var DB *gorm.DB

// Init loads environment variables, connects to the database, and runs migrations.
func Init() {
    loadEnv()
    connectDB()
    migrateDB()
}

// loadEnv reads .env file into environment, if present
func loadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found; relying on environment variables")
    }
}

// connectDB opens a MySQL connection using GORM
func connectDB() {
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASS")
    host := os.Getenv("DB_HOST")
    name := os.Getenv("DB_NAME")
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, name)

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    log.Println("Database connection established")
}

// migrateDB runs AutoMigrate on all models
func migrateDB() {
    if err := DB.AutoMigrate(
        &controllers.User{},
        // add future models here, e.g. &controllers.Profile{}, &controllers.Order{},
    ); err != nil {
        log.Fatalf("failed to run database migrations: %v", err)
    }
    log.Println("Database migrations completed")
}