package initializers

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/leventeberry/goapi/cache"
	"github.com/leventeberry/goapi/config"
	"github.com/leventeberry/goapi/logger"
	"github.com/leventeberry/goapi/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// RedisClient is the global Redis client connection
var RedisClient *redis.Client

// Init loads environment variables, connects to the database, and runs migrations.
func Init() {
	loadEnv()
	validateEnv()
	// Load centralized configuration (must be after loadEnv)
	config.Load()
	connectDB()
	migrateDB()
	connectRedis()
}

// loadEnv reads .env file into environment, if present
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Info().Msg("No .env file found; relying on environment variables")
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
		logger.Log.Fatal().
			Str("missing_vars", fmt.Sprint(missing)).
			Msg("Missing required environment variables. Please create a .env file in the root directory with these variables, or set them in your environment.")
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
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	logger.Log.Info().Msg("Database connection established")
}

// migrateDB runs AutoMigrate on all models
func migrateDB() {
	if err := DB.AutoMigrate(
		&models.User{},
		// add future models here, e.g. &controllers.Profile{}, &controllers.Order{},
	); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	logger.Log.Info().Msg("Database migrations completed")
}

// connectRedis opens a Redis connection if Redis is enabled
// Redis configuration is optional - if REDIS_ENABLED is not "true", Redis will not be connected
func connectRedis() {
	enabled := os.Getenv("REDIS_ENABLED")
	if enabled != "true" {
		logger.Log.Info().Msg("Redis is disabled (REDIS_ENABLED != true)")
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
		logger.Log.Warn().Err(err).Str("address", addr).Msg("Failed to connect to Redis. Application will continue without Redis caching")
		RedisClient = nil
		return
	}

	logger.Log.Info().Str("address", addr).Msg("Redis connection established")
}

// GetCacheClient returns a cache client instance
// Returns Redis cache if Redis is available, otherwise returns no-op cache
// This centralizes cache client creation logic
func GetCacheClient() cache.Cache {
	if RedisClient != nil {
		return cache.NewRedisCache(RedisClient)
	}
	return cache.NewNoOpCache()
}

// CloseRedis closes the Redis connection if it exists
// Should be called on application shutdown for graceful cleanup
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Error closing Redis connection")
		} else {
			logger.Log.Info().Msg("Redis connection closed")
		}
	}
}
