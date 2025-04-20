package middleware

import (
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

const (
    // TokenExpirationDays defines how many days the JWT is valid for.
    TokenExpirationDays = 60
)

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
        jwtSecret := []byte(os.Getenv("JWT_SECRET"))

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
    jwtSecret := []byte(os.Getenv("JWT_SECRET"))
    apiKey := uuid.NewString()
    expiresAt := time.Now().Add(time.Hour * 24 * TokenExpirationDays)

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