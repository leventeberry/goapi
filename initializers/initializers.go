package initializers

import (
    "fmt"
    "log"
    "os"
    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/leventeberry/goapi/models"
)

// DB is the global database connection
var DB *gorm.DB

// Init loads environment variables, connects to the database, and runs migrations.
func Init() {
    loadEnv()
    validateEnv()
    connectDB()
    migrateDB()
}

// loadEnv reads .env file into environment, if present
func loadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found; relying on environment variables")
    }
}

// validateEnv checks that all required environment variables are set
func validateEnv() {
    requiredVars := map[string]string{
        "DB_USER":    "Database username",
        "DB_PASS":    "Database password",
        "DB_HOST":    "Database host (e.g., localhost:3306)",
        "DB_NAME":    "Database name",
        "JWT_SECRET": "JWT secret key for token signing",
    }

    var missing []string
    for key, description := range requiredVars {
        value := os.Getenv(key)
        if value == "" {
            missing = append(missing, fmt.Sprintf("  %s (%s)", key, description))
        }
    }

    if len(missing) > 0 {
        log.Fatalf(
            "Missing required environment variables:\n%s\n\n"+
                "Please create a .env file in the root directory with these variables, or set them in your environment.\n"+
                "Example .env file:\n"+
                "DB_USER=your_db_user\n"+
                "DB_PASS=your_db_password\n"+
                "DB_HOST=localhost:3306\n"+
                "DB_NAME=your_database_name\n"+
                "JWT_SECRET=your_super_secret_jwt_key_here\n"+
                "PORT=8080\n",
            fmt.Sprint(missing),
        )
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
        &models.User{},
        // add future models here, e.g. &controllers.Profile{}, &controllers.Order{},
    ); err != nil {
        log.Fatalf("failed to run database migrations: %v", err)
    }
    log.Println("Database migrations completed")
}