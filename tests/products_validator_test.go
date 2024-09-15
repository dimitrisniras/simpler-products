package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"simpler-products/models"
	"simpler-products/validators"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateProductID(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("ValidID", func(t *testing.T) {
		// Create a request with a valid ID
		req, _ := http.NewRequest("GET", "/products/some-uuid", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "some-uuid"}}

		// Call the validator function
		id, err := validators.ValidateProductID(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, "some-uuid", id)
	})

	t.Run("EmptyID", func(t *testing.T) {
		// Create a request with an empty ID
		req, _ := http.NewRequest("GET", "/products/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		// Call the validator function
		_, err := validators.ValidateProductID(c)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "invalid product id", err.Error())
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})
}

func TestValidateProduct(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("ValidProduct", func(t *testing.T) {
		// Create a request with valid product data
		productData := &models.Product{
			Name:        "Valid Product",
			Description: "This is a valid product",
			Price:       10.99,
		}
		jsonData, _ := json.Marshal(productData)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the validator function
		product, err := validators.ValidateProduct(c)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, productData.Name, product.Name)
		assert.Equal(t, productData.Description, product.Description)
		assert.Equal(t, productData.Price, product.Price)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		// Create a request with missing required fields
		invalidProductData := map[string]interface{}{
			"price": 10.99,
		}
		jsonData, _ := json.Marshal(invalidProductData)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the validator function
		_, err := validators.ValidateProduct(c)

		// Assertions
		assert.Error(t, err)
		var validationErr *validators.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Len(t, validationErr.Errors, 2)
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	t.Run("InvalidPrice", func(t *testing.T) {
		// Create a request with an invalid price (less than or equal to 0)
		invalidProductData := map[string]interface{}{
			"name":        "Invalid Product",
			"description": "This product has an invalid price",
			"price":       0,
		}
		jsonData, _ := json.Marshal(invalidProductData)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the validator function
		_, err := validators.ValidateProduct(c)

		// Assertions
		assert.Error(t, err)
		var validationErr *validators.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Len(t, validationErr.Errors, 1)
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})
}
