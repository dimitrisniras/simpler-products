package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"simpler-products/middlewares"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	// Set up a secret key for JWT (replace with your actual secret key)
	secretKey := "your_secret_key"
	os.Setenv("JWT_SECRET_KEY", secretKey)

	t.Run("ValidToken", func(t *testing.T) {
		// Enable auth
		os.Setenv("AUTH_ENABLED", "true")

		// Create a valid JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
		})
		tokenString, _ := token.SignedString([]byte(secretKey))

		// Create a request with the Authorization header
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		// Assertions
		assert.Equal(t, http.StatusOK, c.Writer.Status())
	})

	t.Run("InvalidToken", func(t *testing.T) {
		// Enable auth
		os.Setenv("AUTH_ENABLED", "true")

		// Create an invalid token (e.g., with a different secret key)
		invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		invalidTokenString, _ := invalidToken.SignedString([]byte("wrong_secret"))

		// Create a request with the Authorization header
		req, _ := http.NewRequest("GET", "/protected", nil)
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

		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, errors.New("invalid token"), err)
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		// Enable auth
		os.Setenv("AUTH_ENABLED", "true")

		// Create a request without the Authorization header
		req, _ := http.NewRequest("GET", "/protected", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, errors.New("authorization header missing"), err)
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("InvalidAuthorizationHeaderFormat", func(t *testing.T) {
		// Enable auth
		os.Setenv("AUTH_ENABLED", "true")

		// Create a request with an invalid Authorization header format
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, errors.New("invalid Authorization header format"), err)
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		// Enable auth
		os.Setenv("AUTH_ENABLED", "true")

		// Create an expired JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(-time.Hour).Unix(), // Expired an hour ago
		})
		tokenString, _ := token.SignedString([]byte(secretKey))

		// Create a request with the Authorization header
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, errors.New("invalid token"), err)
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("AuthDisabled", func(t *testing.T) {
		// Disable auth
		os.Setenv("AUTH_ENABLED", "false")

		// Create a request (Authorization header doesn't matter in this case)
		req, _ := http.NewRequest("GET", "/protected", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the middleware
		middlewares.JWTAuthMiddleware()(c)

		// Assertions
		assert.Equal(t, http.StatusOK, c.Writer.Status())
	})
}
