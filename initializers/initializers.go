package initializers

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/leventeberry/goapi/models"
    "github.com/redis/go-redis/v9"
)

// DB is the global database connection
var DB *gorm.DB

// RedisClient is the global Redis client connection
var RedisClient *redis.Client

// Init loads environment variables, connects to the database, and runs migrations.
func Init() {
    loadEnv()
    validateEnv()
    connectDB()
    migrateDB()
    connectRedis()
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
        "DB_HOST":    "Database host (e.g., localhost)",
        "DB_PORT":    "Database port (e.g., 5432)",
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
                "DB_HOST=localhost\n"+
                "DB_PORT=5432\n"+
                "DB_NAME=your_database_name\n"+
                "JWT_SECRET=your_super_secret_jwt_key_here\n"+
                "PORT=8080\n",
            fmt.Sprint(missing),
        )
    }
}

// connectDB opens a PostgreSQL connection using GORM
func connectDB() {
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASS")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    if port == "" {
        port = "5432"
    }
    name := os.Getenv("DB_NAME")
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, pass, name, port)

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

// connectRedis opens a Redis connection if Redis is enabled
// Redis configuration is optional - if REDIS_ENABLED is not "true", Redis will not be connected
func connectRedis() {
    enabled := os.Getenv("REDIS_ENABLED")
    if enabled != "true" {
        log.Println("Redis is disabled (REDIS_ENABLED != true)")
        return
    }

    host := os.Getenv("REDIS_HOST")
    if host == "" {
        host = "localhost"
    }

    port := os.Getenv("REDIS_PORT")
    if port == "" {
        port = "6379"
    }

    password := os.Getenv("REDIS_PASSWORD")
    // Password is optional, empty string means no password

    addr := fmt.Sprintf("%s:%s", host, port)
    
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0, // Default DB
    })

    // Test connection
    ctx := context.Background()
    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Printf("Warning: Failed to connect to Redis at %s: %v", addr, err)
        log.Println("Application will continue without Redis caching")
        RedisClient = nil
        return
    }

    log.Printf("Redis connection established at %s", addr)
}