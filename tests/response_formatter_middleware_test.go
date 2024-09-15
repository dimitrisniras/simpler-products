// tests/response_formatter_test.go (create this file in your tests directory)
package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"simpler-products/middlewares"
	"simpler-products/validators"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestResponseFormatter(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	// Create a mock logger
	log := logrus.New()

	t.Run("SuccessWithData", func(t *testing.T) {
		// Create a Gin context and set data
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("data", []string{"item1", "item2"})

		// Call the middleware
		middlewares.ResponseFormatter(log)(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, int(response["status"].(float64)))
		assert.Equal(t, []interface{}{"item1", "item2"}, response["data"])
		assert.Nil(t, response["errors"])
		assert.Nil(t, response["pagination"])
	})

	t.Run("SuccessWithPagination", func(t *testing.T) {
		// Create a Gin context and set data and pagination
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("data", []string{"item1", "item2"})
		c.Set("pagination", gin.H{
			"limit":  10,
			"offset": 0,
			"total":  25,
			"count":  10,
		})

		// Call the middleware
		middlewares.ResponseFormatter(log)(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, int(response["status"].(float64)))
		assert.Equal(t, []interface{}{"item1", "item2"}, response["data"])
		assert.Equal(t, map[string]interface{}{"limit": 10.0, "offset": 0.0, "total": 25.0, "count": 10.0}, response["pagination"])
		assert.Nil(t, response["errors"])
	})

	t.Run("ValidationError", func(t *testing.T) {
		// Create a Gin context and set a validation error
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("errors", &validators.ValidationError{
			Errors: []map[string]string{
				{"message": "Name is required"},
				{"message": "Price must be greater than 0"},
			},
		})
		c.Writer.WriteHeader(http.StatusBadRequest) // Set the status code

		// Call the middleware
		middlewares.ResponseFormatter(log)(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, int(response["status"].(float64)))
		assert.Nil(t, response["data"])
		assert.Nil(t, response["pagination"])
		assert.Equal(t, []interface{}{
			map[string]interface{}{"message": "Name is required"},
			map[string]interface{}{"message": "Price must be greater than 0"},
		}, response["errors"])
	})

	t.Run("OtherError", func(t *testing.T) {
		// Create a Gin context and set a generic error
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("errors", errors.New("some error occurred"))
		c.Writer.WriteHeader(http.StatusInternalServerError)

		// Call the middleware
		middlewares.ResponseFormatter(log)(c)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, int(response["status"].(float64)))
		assert.Nil(t, response["data"])
		assert.Nil(t, response["pagination"])
		assert.Equal(t, []interface{}{map[string]interface{}{"message": "some error occurred"}}, response["errors"])
	})

	t.Run("NoDataNoErrors", func(t *testing.T) {
		// Create a Gin context without setting any data or errors
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the middleware
		middlewares.ResponseFormatter(log)(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, int(response["status"].(float64)))
		assert.Nil(t, response["data"])
		assert.Nil(t, response["errors"])
		assert.Nil(t, response["pagination"])
	})

	t.Run("DifferentStatusCodes", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
			data       interface{}
			errors     interface{}
			expected   map[string]interface{}
		}{
			{
				name:       "201Created",
				statusCode: http.StatusCreated,
				data:       map[string]string{"message": "created"},
				expected: map[string]interface{}{
					"status": float64(http.StatusCreated),
					"data":   map[string]interface{}{"message": "created"},
				},
			},
			{
				name:       "400BadRequest",
				statusCode: http.StatusBadRequest,
				errors:     []string{"Invalid request"},
				expected: map[string]interface{}{
					"status": float64(http.StatusBadRequest),
					"errors": []interface{}{map[string]interface{}{"message": "Invalid request"}},
				},
			},
			{
				name:       "401Unauthorized",
				statusCode: http.StatusUnauthorized,
				errors:     []string{"Unauthorized request"},
				expected: map[string]interface{}{
					"status": float64(http.StatusUnauthorized),
					"errors": []interface{}{map[string]interface{}{"message": "Unauthorized request"}},
				},
			},
			{
				name:       "404NotFound",
				statusCode: http.StatusNotFound,
				errors:     errors.New("resource not found"),
				expected: map[string]interface{}{
					"status": float64(http.StatusNotFound),
					"errors": []interface{}{map[string]interface{}{"message": "resource not found"}},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create a Gin context and set the status code and data/errors
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				if tc.data != nil {
					c.Set("data", tc.data)
				}
				if tc.errors != nil {
					c.Set("errors", tc.errors)
				}
				c.Writer.WriteHeader(tc.statusCode)

				// Call the middleware
				middlewares.ResponseFormatter(log)(c)

				// Assertions
				assert.Equal(t, tc.statusCode, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tc.expected, response)
			})
		}
	})
}
