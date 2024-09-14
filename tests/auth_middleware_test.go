// tests/auth_middleware_test.go
package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	custom_errors "simpler-products/errors"
	"simpler-products/middlewares"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	// Generate RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Error generating private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Encode public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("Error marshalling public key: %v", err)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)

	// Encode public key to Base64 and set as environment variable
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyPEM)
	os.Setenv("JWT_SECRET_KEY", publicKeyBase64)

	t.Run("ValidToken", func(t *testing.T) {
		// Create a valid JWT token using the private key
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(privateKey)
		assert.NoError(t, err)

		// Create a request with the Authorization header
		req, _ := http.NewRequest("GET", "/api/v1/v1/products", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set AUTH_ENABLED to true
		os.Setenv("AUTH_ENABLED", "true")

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		// Assertions
		assert.Equal(t, http.StatusOK, c.Writer.Status()) // Should proceed to the next handler
	})

	t.Run("InvalidToken", func(t *testing.T) {
		// Create an invalid token (e.g., with a different private key)
		invalidPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("Error generating invalid private key: %v", err)
		}

		invalidToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		invalidTokenString, _ := invalidToken.SignedString(invalidPrivateKey)

		// Create a request with the Authorization header
		req, _ := http.NewRequest("GET", "/api/v1/products", nil)
		req.Header.Set("Authorization", "Bearer "+invalidTokenString)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set AUTH_ENABLED to true
		os.Setenv("AUTH_ENABLED", "true")

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err_, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
		assert.Equal(t, custom_errors.ErrInvalidToken, err_)
	})

	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		// Create a request without the Authorization header
		req, _ := http.NewRequest("GET", "/api/v1/products", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set AUTH_ENABLED to true
		os.Setenv("AUTH_ENABLED", "true")

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err_, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
		assert.Equal(t, custom_errors.ErrAuthorizationHeaderMissing, err_)
	})

	t.Run("InvalidAuthorizationHeaderFormat", func(t *testing.T) {
		// Create a request with an invalid Authorization header format
		req, _ := http.NewRequest("GET", "/api/v1/products", nil)
		req.Header.Set("Authorization", "InvalidFormat")

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set AUTH_ENABLED to true
		os.Setenv("AUTH_ENABLED", "true")

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err_, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
		assert.Equal(t, custom_errors.ErrAuthorizationHeaderFormat, err_)
	})

	t.Run("ErrorDecodingPublicKey", func(t *testing.T) {
		// Set an invalid Base64 string as the public key
		os.Setenv("JWT_SECRET_KEY", "invalid_base64")

		// Create a request with a valid token (but it won't be verified due to the invalid public key)
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, _ := token.SignedString(privateKey)
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set AUTH_ENABLED to true
		os.Setenv("AUTH_ENABLED", "true")

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err_, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
		assert.Equal(t, custom_errors.ErrDecodingPublicKey, err_)
	})
}
