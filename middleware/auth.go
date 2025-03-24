package middleware

import (
	"net/http"
	"strings"
	"time"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RequestTokens struct {
	ApiKey string
	JWT_Token string
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("exp", claims["exp"])
		}

		c.Next()
	}
}


func CreateToken(userID int) (RequestTokens, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	uuid := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24 * 60).Unix() // 60 Days

	claims := jwt.MapClaims{
		"apiKey":  uuid,
		"exp":     expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return RequestTokens{}, err
	}

	return RequestTokens{
		ApiKey: uuid,
		JWT_Token: signedToken,
	}, nil
}
