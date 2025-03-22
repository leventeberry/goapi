package middleware

import (
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your-secret-key") // store securely, e.g. in .env or config

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Optionally set user info in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("userID", claims["user_id"])
			c.Set("email", claims["email"])
			// Add more claims if needed
		}

		c.Next() // Allow request to continue
	}
}

func CreateToken(userID int) (string, error) {
	uuid := uuid.New()
	
	// Token expires in 1 hour
	expiresAt := time.Now().Add(time.Hour * 60).Unix()

	claims := jwt.MapClaims{
		"apiKey":  uuid.String(),
		"exp":     expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
