package middlewares

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	custom_errors "simpler-products/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authEnabled := os.Getenv("AUTH_ENABLED")

		if authEnabled == "true" {
			// Get the Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", custom_errors.ErrAuthorizationHeaderMissing)
				return
			}

			// Split the header to get the token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", custom_errors.ErrAuthorizationHeaderFormat)
				return
			}
			tokenString := parts[1]

			secretKey := os.Getenv("JWT_SECRET_KEY")

			// Decode the Base64-encoded public key
			publicKeyBytes, err := base64.StdEncoding.DecodeString(secretKey)
			if err != nil {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", custom_errors.ErrDecodingPublicKey)
				return
			}

			// Parse the public key
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
			if err != nil {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", custom_errors.ErrParsingPublicKey)
				return
			}

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Make sure the signing method is RSA
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return publicKey, nil
			})

			if err != nil || !token.Valid {
				c.Status(http.StatusUnauthorized)
				c.Set("errors", custom_errors.ErrInvalidToken)
				return
			}
		}

		// Token is valid, continue to the next handler
		c.Next()
	}
}
