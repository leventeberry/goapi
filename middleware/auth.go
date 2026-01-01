package middleware

import (
    "net/http"
    "strconv"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/leventeberry/goapi/config"
    "github.com/leventeberry/goapi/repositories"
)

// getTokenExpirationDays returns the JWT token expiration days from configuration
func getTokenExpirationDays() int {
	cfg := config.Get()
	return cfg.JWT.ExpirationDays
}

// Claims defines the JWT payload structure.
type Claims struct {
    ApiKey string `json:"api_key"`
    jwt.RegisteredClaims
}

// TokenDetails holds the generated API key and JWT token.
type Authentication struct {
    ApiKey   string `json:"api_key"`
    JWTToken string `json:"jwt_token"`
}

// AuthMiddleware validates the JWT token from the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        cfg := config.Get()
        jwtSecret := []byte(cfg.JWT.Secret)

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return jwtSecret, nil
        })
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }

        // Store claims in context
        c.Set("apiKey", claims.ApiKey)
        c.Set("userID", claims.Subject)
        c.Set("expiresAt", claims.ExpiresAt.Time)

        c.Next()
    }
}

// CreateToken generates a new JWT token (and API key) for the given user ID.
func CreateToken(userID int) (*Authentication, error) {
    cfg := config.Get()
    jwtSecret := []byte(cfg.JWT.Secret)
    apiKey := uuid.NewString()
    expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(getTokenExpirationDays()))

    claims := Claims{
        ApiKey: apiKey,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   strconv.Itoa(userID),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return nil, err
    }

    return &Authentication{
        ApiKey:   apiKey,
        JWTToken: signedToken,
    }, nil
}

// RequireRole returns a middleware that checks if the authenticated user has one of the required roles.
// This middleware must be used after AuthMiddleware, as it relies on userID being set in the context.
// Uses dependency injection to access user repository
func RequireRole(userRepo repositories.UserRepository, allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get userID from context (set by AuthMiddleware)
        userIDStr, exists := c.Get("userID")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
            return
        }

        // Convert userID string to int
        userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
        if err != nil || userID < 1 {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
            return
        }

        // Query database for user's role using repository
        user, err := userRepo.FindByID(int(userID))
        if err != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        // Check if user's role is in the allowed roles list
        hasRole := false
        for _, role := range allowedRoles {
            if user.Role == role {
                hasRole = true
                break
            }
        }

        if !hasRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            return
        }

        // Store user object in context for use in handlers
        c.Set("user", user)
        c.Next()
    }
}