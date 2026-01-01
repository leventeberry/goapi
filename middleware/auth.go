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
)

// getTokenExpirationDays returns the JWT token expiration days from configuration
func getTokenExpirationDays() int {
	cfg := config.Get()
	return cfg.JWT.ExpirationDays
}

// Claims defines the JWT payload structure.
type Claims struct {
    ApiKey string `json:"api_key"`
    Role   string `json:"role"`
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
        c.Set("role", claims.Role)
        c.Set("expiresAt", claims.ExpiresAt.Time)

        c.Next()
    }
}

// CreateToken generates a new JWT token (and API key) for the given user ID and role.
func CreateToken(userID int, role string) (*Authentication, error) {
    cfg := config.Get()
    jwtSecret := []byte(cfg.JWT.Secret)
    apiKey := uuid.NewString()
    expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(getTokenExpirationDays()))

    claims := Claims{
        ApiKey: apiKey,
        Role:   role,
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
// This middleware must be used after AuthMiddleware, as it relies on role being set in the context.
// Role is now stored in JWT token claims, eliminating the need for database queries.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get role from context (set by AuthMiddleware from JWT claims)
        role, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
            return
        }

        roleStr, ok := role.(string)
        if !ok || roleStr == "" {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format in token"})
            return
        }

        // Check if user's role is in the allowed roles list
        hasRole := false
        for _, allowedRole := range allowedRoles {
            if roleStr == allowedRole {
                hasRole = true
                break
            }
        }

        if !hasRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            return
        }

        c.Next()
    }
}