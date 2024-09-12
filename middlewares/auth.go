package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET_KEY")
		authEnabled := os.Getenv("AUTH_ENABLED")

		if boolValue, err := strconv.ParseBool(authEnabled); err == nil && boolValue == true {
			// Get the Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", errors.New("authorization header missing"))
				return
			}

			// Split the header to get the token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", errors.New("invalid Authorization header format"))
				return
			}
			tokenString := parts[1]

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Make sure the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", errors.New("invalid token"))
				return

			}
		}

		// Token is valid, continue to the next handler
		c.Next()
	}
}
